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
	expireDuration time.Duration
}

func NewAuthUseCase(userRepo auth.UserRepository,roleRepo role.RoleRepository, hashSalt string, signingKey []byte, tokenTTL int64) auth.UseCase {
	return &authUserCase{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * time.Duration(tokenTTL),
	}
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

func (a *authUserCase) SignIn(ctx context.Context, creds presenter.LoginInput) (string, error) {
	user, _ := a.userRepo.GetByField(ctx, "Username", creds.Username)
	if user == nil {
		return "", auth.ErrUserNotFound
	}
	if !user.ComparePassword(creds.Password) {
		return "", auth.ErrWrongPassword
	}
	var idsString = make([]string, len(user.RoleIds))
	for i,id := range user.RoleIds {
		idsString[i] = id.Hex()
	}
	roles, err := a.roleRepo.GetMany(ctx, idsString)
	if err != nil {
		return "", err
	}
	var rolesName = make([]string, len(*roles))
	for i, role := range *roles {
		rolesName[i] = role.Name
	}
	claims := AuthClaims{
		Username: user.Username,
		UserId:   user.Id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    "saving-books",
			ExpiresAt: time.Now().Add(a.expireDuration).Unix(),
		},
		Roles: rolesName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.signingKey)
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
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return &presenter.TokenResult{
			UserId: claims.UserId.Hex(),
			Roles: claims.Roles,
		}, nil
	}
	return nil, auth.ErrInvalidAccessToken
}
