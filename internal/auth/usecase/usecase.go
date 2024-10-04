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
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthClaims struct {
	jwt.StandardClaims
	Username string             `json:"username"`
	UserId   primitive.ObjectID `json:"userId"`
}
type authUserCase struct {
	userRepo             auth.UserRepository
	hashSalt             string
	signingKey           []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func (a *authUserCase) Logout(ctx context.Context, userId string) error {
	user, err := a.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return err
	}
	user.AssignRefreshToken(uuid.New().String(), time.Now().Add(a.refreshTokenDuration))
	err = a.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (a *authUserCase) RenewAccessToken(ctx context.Context, req *presenter.RenewTokenReq) (string, error) {
	user, err := a.userRepo.GetUserById(ctx, req.UserId)
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
	euser, _ := a.userRepo.GetUserByUsername(ctx, fmtusersame)
	if euser != nil {
		return nil, auth.ErrUserExisted
	}
	aggregateRoot := &domain.AggregateRoot{}
	aggregateRoot.SetInit()
	user := &domain.User{
		AggregateRoot: *aggregateRoot,
		Username:      creds.Username,
		Password:      creds.Password,
	}
	user.HashPassword()
	err := a.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (a *authUserCase) SignIn(ctx context.Context, creds presenter.LoginInput) (*presenter.LogInRes, error) {
	user, _ := a.userRepo.GetUserByUsername(ctx, creds.Username)
	if user == nil {
		return nil, auth.ErrUserNotFound
	}
	if !user.ComparePassword(creds.Password) {
		return nil, auth.ErrWrongPassword
	}
	refreshToken := uuid.New().String()
	user.AssignRefreshToken(refreshToken, time.Now().Add(a.refreshTokenDuration))
	err := a.userRepo.UpdateUser(ctx, user)
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

func (a *authUserCase) ParseAccessToken(accessToken string) (*presenter.ParseTokenResult, error) {
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
		return &presenter.ParseTokenResult{
			UserId: claims.UserId.Hex(),
		}, nil
	}
	return nil, auth.ErrInvalidAccessToken
}
func NewAuthUseCase(userRepo auth.UserRepository, hashSalt string, signingKey []byte, accessTokenDuration int64, refreshTokenDuration int64) auth.UseCase {
	return &authUserCase{
		userRepo:             userRepo,
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
