package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/domain"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthClaims struct {
	jwt.StandardClaims
	Username string             `json:"username"`
	UserId   primitive.ObjectID `json:"userId"`
}
type authUserCase struct {
	userRepo       auth.UserRepository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(userRepo auth.UserRepository, hashSalt string, signingKey []byte, tokenTTL int64) auth.UseCase {
	return &authUserCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * time.Duration(tokenTTL),
	}
}
func (a *authUserCase) SignUp(ctx context.Context, creds presenter.SignUpInput) (*domain.User, error) {
	fmtusersame := strings.ToLower(creds.Username)
	euser, _ := a.userRepo.GetUserByUsername(ctx, fmtusersame)
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
	err := a.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authUserCase) SignIn(ctx context.Context, creds presenter.LoginInput) (string, error) {
	user, _ := a.userRepo.GetUserByUsername(ctx, creds.Username)
	if user == nil {
		return "", auth.ErrUserNotFound
	}
	if !user.ComparePassword(creds.Password) {
		return "", auth.ErrWrongPassword
	}
	claims := AuthClaims{
		Username: user.Username,
		UserId:   user.Id,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			Issuer:    "saving-books",
			ExpiresAt: time.Now().Add(a.expireDuration).Unix(),
		},
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
			Username: claims.Username,
		}, nil
	}
	return nil, auth.ErrInvalidAccessToken
}
