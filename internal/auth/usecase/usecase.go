package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthClaims struct {
	jwt.StandardClaims
	Username string             `json:"username"`
	UserId   primitive.ObjectID `json:"userId"`
	Roles    []string           `json:"roles"`

}
type authUserCase struct {
	userRepo       auth.UserRepository
	roleRepo       role.RoleRepository
	hashSalt       string
	signingKey     []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func (a *authUserCase) Logout(ctx context.Context, userId string) error {
	user, err := a.userRepo.Get(ctx, userId)
	if err != nil {
		return err
	}
	user.AssignRefreshToken(uuid.New().String(), time.Now().Add(a.refreshTokenDuration))
	err = a.userRepo.Update(ctx, user, user.Id.Hex())
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
	if user.RefreshToken != req.RefreshToken {
		return "", auth.ErrInvalidRefreshToken
	}
	if user.RefreshTokenExpiresAt.Before(time.Now()) {
		return "", auth.ErrRefreshTokenExpired
	}
	token, err := a.generateAccessToken(user)
	if err != nil {
		return "", err

	}
	return token, err
}
func (a *authUserCase) SignUp(ctx context.Context, creds presenter.SignUpInput) (*domain.User, error) {
	fmtusersame := strings.ToLower(creds.Username)
	euser, _ := a.userRepo.GetByField(ctx, "Username", fmtusersame)
	if euser != nil {
		return nil, auth.ErrUserExisted
	}
	aggregateRoot := &domain.AggregateRoot{}
	aggregateRoot.SetInit()
	user := &domain.User{
		AggregateRoot: *aggregateRoot,
		Username: creds.Username,
		Password: creds.Password,
	}
	user.HashPassword()

	count, err := a.userRepo.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		roleT, ierr := a.roleRepo.GetByField(ctx ,"Name", "Admin")
		if ierr != nil {
			return nil, ierr
		}
		user.RoleIds =  append(user.RoleIds, roleT.Id)
	}else {
		roleT, ierr := a.roleRepo.GetByField(ctx ,"Name", "User")
		if ierr != nil {
			return nil, ierr
		}
		user.RoleIds =  append(user.RoleIds, roleT.Id)
	}
	err = a.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authUserCase) SignIn(ctx context.Context, creds presenter.LoginInput) (*presenter.LogInRes, error) {
	user, _ := a.userRepo.GetByField(ctx, "Username", creds.Username)
	if user == nil {
		return nil, auth.ErrUserNotFound
	}
	if !user.ComparePassword(creds.Password) {
		return nil, auth.ErrWrongPassword
	}
	var idsString = make([]string, len(user.RoleIds))
	for i,id := range user.RoleIds {
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
	refreshToken := uuid.New().String()
	user.AssignRefreshToken(refreshToken, time.Now().Add(a.refreshTokenDuration))
	err = a.userRepo.Update(ctx, user, user.Id.Hex())
	if err != nil {
		return nil, err
	}
	tokenString, err := a.generateAccessToken(user)

	if err != nil {
		return nil, err
	}

	return &presenter.LogInRes{
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authUserCase) ParseAccessToken(accessToken string) (*presenter.TokenResult, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return &presenter.TokenResult{
			UserId: claims.UserId.Hex(),
			Roles: claims.Roles,
		}, nil
	}
	return nil, auth.ErrInvalidAccessToken
}
func NewAuthUseCase(userRepo auth.UserRepository, roleRepo role.RoleRepository, hashSalt string, signingKey []byte, accessTokenDuration int64, refreshTokenDuration int64) auth.UseCase {
	return &authUserCase{
		userRepo:             userRepo,
		roleRepo:             roleRepo,
		hashSalt:             hashSalt,
		signingKey:           signingKey,
		accessTokenDuration:  time.Second * time.Duration(accessTokenDuration),
		refreshTokenDuration: time.Second * time.Duration(refreshTokenDuration),
	}
}

func (a *authUserCase) generateAccessToken(user *domain.User) (string, error) {
	accessClaims := AuthClaims{
		Username: user.Username,
		UserId:   user.Id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    "saving-books",
			ExpiresAt: time.Now().Add(a.accessTokenDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	tokenString, err := token.SignedString(a.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, err
}
