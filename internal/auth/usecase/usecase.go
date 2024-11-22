package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"SavingBooks/internal/services/email"
	"SavingBooks/internal/services/redis"
	"SavingBooks/internal/services/redis/redis_key"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthClaims struct {
	jwt.StandardClaims
	Username string             `json:"username"`
	UserId   primitive.ObjectID `json:"userId"`
	Roles    []string           `json:"roles"`
}
type authUserCase struct {
	userRepo             auth.UserRepository
	roleRepo             role.RoleRepository
	emailService         *email.SmtpServer
	cacheService *redis.Cache
	hashSalt             string
	signingKey           []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func (a *authUserCase) ChangePassword(ctx context.Context,userId, oldPassword, newPassword string) error {

	user, err := a.userRepo.Get(ctx, userId)
	if err != nil {
		return err
	}

	if !user.ComparePassword(oldPassword) {
		return auth.ErrWrongPassword
	}
	user.Password = newPassword
	user.HashPassword()
	user.RefreshToken = ""
	_, err = a.userRepo.Update(ctx, user, userId, []string{"Password", "RefreshToken"})
	if err != nil {
		return err
	}
	return nil
}

func (a *authUserCase) GenerateResetPassword(ctx context.Context, email string) error {
	user, err := a.userRepo.GetByField(ctx, "Email", email)
	if err != nil {
		return err
	}
	if !user.IsActive {
		return auth.ErrUserIsBlocked
	}
	user.GenerateResetPassword()
	_, err = a.userRepo.Update(ctx, user, user.Id.Hex(), []string{"ResetPasswordToken", "ResetPasswordExpiresAt" })
	if err != nil {
		return err
	}
	resetLink := fmt.Sprintf("http://localhost:3003/confirm-reset-password?token=%s", user.ResetPasswordToken)
	err = a.emailService.SendEmail(user.Email, resetLink)
	if err != nil {
		return err
	}
	return nil
}
func (a *authUserCase) ConfirmResetPassword(ctx context.Context, token, password string) error{
	user, err := a.userRepo.GetByField(ctx, "ResetPasswordToken", token)
	if err != nil {
		return auth.ErrInvalidResetPasswordToken
	}
	if !user.IsActive {
		return auth.ErrUserIsBlocked
	}
	if user.ResetPasswordToken != token {
		return auth.ErrInvalidResetPasswordToken
	}
	test := time.Now()
	if test.After(user.ResetPasswordExpiresAt) {
		return auth.ErrInvalidResetPasswordToken
	}
	user.Password = password
	user.ResetPasswordToken = ""
	user.HashPassword()

	_, err = a.userRepo.Update(ctx, user, user.Id.Hex(), []string{"Password", "ResetPasswordToken" })
	if err != nil {
		return err
	}
	return nil
}

func (a *authUserCase) Logout(ctx context.Context, userId string) error {
	user, err := a.userRepo.Get(ctx, userId)
	if err != nil {
		return err
	}
	user.AssignRefreshToken(uuid.New().String(), time.Now().Add(a.refreshTokenDuration))
	_, err = a.userRepo.Update(ctx, user, user.Id.Hex(), nil)
	if err != nil {
		return err
	}
	return nil
}

func (a *authUserCase) RenewAccessToken(ctx context.Context, req *presenter.RenewTokenReq) (string, error) {
	user, err := a.userRepo.Get(ctx, req.UserId)
	if err != nil {
		return "", err
	}
	if !user.IsActive {
		return "", auth.ErrUserIsBlocked
	}
	if user.RefreshToken != req.RefreshToken {
		return "", auth.ErrInvalidRefreshToken
	}
	if user.RefreshTokenExpiresAt.Before(time.Now()) {
		return "", auth.ErrRefreshTokenExpired
	}
	rolesName, err := a.getRolesName(ctx, user.RoleIds)
	token, err := a.generateAccessToken(user, rolesName)
	if err != nil {
		return "", err

	}
	return token, err
}
func (a *authUserCase) SignUp(ctx context.Context, creds presenter.SignUpInput) (*domain.User, error) {
	fmtUsername := strings.ToLower(creds.Username)
	fmtEmail := strings.ToLower(creds.Email)
	existUser, err := a.userRepo.GetExistUser(ctx, fmtUsername, fmtEmail)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	}
	if existUser != nil {
		return nil, auth.ErrUserExisted
	}

	aggregateRoot := &domain.AggregateRoot{}
	aggregateRoot.SetInit()
	user := &domain.User{
		AggregateRoot: *aggregateRoot,
		Username:      fmtUsername,
		Password:      creds.Password,
		Email:         fmtEmail,
	}
	user.HashPassword()

	count, err := a.userRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		roleT, ierr := a.roleRepo.GetByField(ctx, "Name", "Admin")
		if ierr != nil {
			return nil, ierr
		}
		user.RoleIds = append(user.RoleIds, roleT.Id)
	} else {
		roleT, ierr := a.roleRepo.GetByField(ctx, "Name", "User")
		if ierr != nil {
			return nil, ierr
		}
		user.RoleIds = append(user.RoleIds, roleT.Id)
	}
	err = a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authUserCase) SignIn(ctx context.Context, creds presenter.LoginInput) (*presenter.LogInRes, error) {
	fmtUsername := strings.ToLower(creds.Username)
	user, err := a.userRepo.GetExistUser(ctx, fmtUsername, fmtUsername)
	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
	}
	if user == nil {
		return nil, auth.ErrLoginCredentials
	}
	if !user.IsActive {
		return nil, auth.ErrUserIsBlocked
	}
	if !user.ComparePassword(creds.Password) {
		return nil, auth.ErrLoginCredentials
	}
	rolesName, err := a.getRolesName(ctx, user.RoleIds)
	refreshToken := uuid.New().String()
	user.AssignRefreshToken(refreshToken, time.Now().Add(a.refreshTokenDuration))
	_, err = a.userRepo.Update(ctx, user, user.Id.Hex(), []string{"RefreshToken", "RefreshTokenExpiresAt"})
	if err != nil {
		return nil, err
	}
	tokenString, err := a.generateAccessToken(user, rolesName)

	if err != nil {
		return nil, err
	}

	return &presenter.LogInRes{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authUserCase) ParseAccessToken(ctx context.Context, accessToken string) (*presenter.TokenResult, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})
	if err != nil {
		return nil, err
	}



	setRoles := map[string]interface{}{}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {

		var userid interface{}
		_ = a.cacheService.GetValue(ctx, redis_key.BlockUserId+":" + claims.UserId.Hex(), &userid)
		if  userid != nil {
			return nil, auth.ErrUserIsBlocked
		}


		for _, role := range claims.Roles {
			setRoles[role] = struct {
			}{}
		}
		return &presenter.TokenResult{
			UserId: claims.UserId.Hex(),
			Roles:  setRoles,
		}, nil
	}
	return nil, auth.ErrInvalidAccessToken
}
func NewAuthUseCase(userRepo auth.UserRepository, roleRepo role.RoleRepository, emailService         *email.SmtpServer,cacheService *redis.Cache, hashSalt string, signingKey []byte, accessTokenDuration int64, refreshTokenDuration int64) auth.UseCase {
	return &authUserCase{
		userRepo:             userRepo,
		roleRepo:             roleRepo,
		emailService:         emailService,
		cacheService: cacheService,
		hashSalt:             hashSalt,
		signingKey:           signingKey,
		accessTokenDuration:  time.Second * time.Duration(accessTokenDuration),
		refreshTokenDuration: time.Second * time.Duration(refreshTokenDuration),
	}
}

func (a *authUserCase) generateAccessToken(user *domain.User, rolesName []string) (string, error) {
	accessClaims := AuthClaims{
		Username: user.Username,
		UserId:   user.Id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    "saving-books",
			ExpiresAt: time.Now().Add(a.accessTokenDuration).Unix(),
		},
		Roles: rolesName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	tokenString, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
func (a *authUserCase) getRolesName(ctx context.Context, roleIds []primitive.ObjectID) ([]string, error) {
	var idsString = make([]string, len(roleIds))
	for i, id := range roleIds {
		idsString[i] = id.Hex()
	}
	roles, err := a.roleRepo.GetMany(ctx, idsString)
	if err != nil {
		return nil, err
	}
	var rolesName = make([]string, len(*roles))
	for i, role := range *roles {
		rolesName[i] = role.Name
	}
	return rolesName, nil
}
