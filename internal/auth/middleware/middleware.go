package middleware

import "SavingBooks/internal/auth"

type MiddleWareManager struct {
	authUC auth.UseCase
}

func NewMiddleWareManager(authUC auth.UseCase) *MiddleWareManager  {
	return &MiddleWareManager{authUC}
}
