package repository

import (
	"context"

	"SavingBooks/internal/auth"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	domain.GenericRepository[domain.User]
}


func (ur *userRepository) GetExistUser(ctx context.Context, username, email string) (*domain.User, error) {
	collectionInterface := ur.GetCollection()
	collection := collectionInterface.(*mongo.Collection)

	filter := bson.M{
		"$or": []bson.M{
			{"Username": username},
			{"Email": email},
		},
	}
	user := &domain.User{}
	err := collection.FindOne(ctx, filter ).Decode(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//collectionInterface := s.GetCollection()
//collection := collectionInterface.(*mongo.Collection)
//
//filter := bson.M{"IsActive": true, "IsDeleted": false}
//reg := &domain.SavingRegulation{}
//opts := options.FindOne().SetSort(bson.D{{"CreationTime", -1}})
//
//err := collection.FindOne(ctx, filter, opts).Decode(&reg)
//if err != nil {
//if errors.Is(err, mongo.ErrNoDocuments) {
//return nil, errors.New(saving_regulation.SavingRegulationNotFoundError)
//}
//return nil, err
//}
//return reg, nil

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
