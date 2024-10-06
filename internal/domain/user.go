package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	AggregateRoot `bson:",inline" json:",inline"`
	Username string `bson:"Username" json:"username"`
	Password string `bson:"Password" json:"password"`
	RoleIds []primitive.ObjectID `bson:"RoleIds" json:"rolesIds"`
	RefreshToken          string    `bson:"RefreshToken" json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `bson:"RefreshTokenExpiresAt" json:"refresh_token_expires_at"`
}

func (u *User) AssignRefreshToken(refreshToken string, expireAt time.Time) {
	u.RefreshToken = refreshToken
	u.RefreshTokenExpiresAt = expireAt
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
func (u *User) ComparePassword(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}
	return true
}
