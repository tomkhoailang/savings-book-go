package auth

import "SavingBooks/internal/domain"

type UserRepository interface {
	domain.GenericRepository[domain.User]
}
