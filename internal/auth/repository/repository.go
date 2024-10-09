package repository

import (
	"SavingBooks/internal/auth"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	domain.GenericRepository[domain.User]
}

func NewUserRepository(db *mongo.Database, collectionName string) auth.UserRepository {
	baseRepo := contracts.NewBaseRepository[domain.User](db, collectionName).(*contracts.BaseRepository[domain.User])
	return &userRepository{GenericRepository: baseRepo}
}

//func (ur *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
//	collection := ur.db.Database(ur.cfg.DatabaseName).Collection(ur.collectionName)
//	_, err := collection.InsertOne(ctx, user)
//	return err
//}
//
//func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
//	collection := ur.db.Database(ur.cfg.DatabaseName).Collection(ur.collectionName)
//
//	var user domain.User
//	err := collection.FindOne(ctx, bson.M{"Username": username}).Decode(&user)
//	if err != nil {
//		if errors.Is(err, mongo.ErrNoDocuments) {
//			return nil, nil
//		}
//		return nil, err
//	}
//	return &user, nil
//}
//
//func (ur *userRepository) GetUserById(ctx context.Context, id string) (*domain.User, error) {
//	collection := ur.db.Database(ur.cfg.DatabaseName).Collection(ur.collectionName)
//
//	var user domain.User
//	err := collection.FindOne(ctx, bson.M{"_id": id, "IsDeleted": false}).Decode(&user)
//	if err != nil {
//		if errors.Is(err, mongo.ErrNoDocuments) {
//			return nil, nil
//		}
//		return nil, err
//	}
//	return &user, nil
//}
