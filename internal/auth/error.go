package auth

import "errors"

var (
	ErrUserNotFound        = errors.New("User not found")
	ErrUserIsBlocked        = errors.New("User is blocked")

	ErrWrongPassword       = errors.New("Old password incorrect")
	ErrLoginCredentials  = errors.New("Login credentials incorrect")
	ErrUserExisted         = errors.New("User already exists")
	ErrInvalidAccessToken  = errors.New("Invalid access token")
	ErrInvalidRefreshToken = errors.New("Invalid refresh token")
	ErrRefreshTokenExpired = errors.New("Refresh token expired")
	ErrInvalidResetPasswordToken = errors.New("Either reset token is invalid or expired")
)
