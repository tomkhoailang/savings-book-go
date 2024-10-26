package auth

import "errors"

var (
	ErrUserNotFound        = errors.New("User not found")
	ErrWrongPassword       = errors.New("Password incorrect")
	ErrLoginCredentials  = errors.New("Login credentials incorrect")
	ErrUserExisted         = errors.New("User already exists")
	ErrInvalidAccessToken  = errors.New("Invalid access token")
	ErrInvalidRefreshToken = errors.New("Invalid refresh token")
	ErrRefreshTokenExpired = errors.New("Refresh token expired")
)
