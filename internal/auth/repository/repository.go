package repository

import (
	"context"
	"errors"

	"SavingBooks/config"
	"SavingBooks/internal/auth"
	"SavingBooks/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db     *mongo.Client
	cfg    *config.Configuration
	collectionName string
}

func NewUserRepository(db *mongo.Client, cfg *config.Configuration) auth.UserRepository {
	return &userRepository{db: db, cfg: cfg, collectionName: "User"}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	collection := ur.db.Database(ur.cfg.DatabaseName).Collection(ur.collectionName)
	_, err := collection.InsertOne(ctx, user)
	return err
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	collection := ur.db.Database(ur.cfg.DatabaseName).Collection(ur.collectionName)

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"Username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	collection := ur.db.Database(ur.cfg.DatabaseName).Collection(ur.collectionName)

	var user domain.User
	err := collection.FindOne(ctx, bson.M{"_id": id , "is_deleted" : false}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
