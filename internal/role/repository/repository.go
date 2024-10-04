package repository

import (
	"context"

	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/role"
	"go.mongodb.org/mongo-driver/mongo"
)

type roleRepository struct {
	contracts.BaseRepository[domain.Role]
}

func (r roleRepository) TestMethod(ctx context.Context, id string) (*domain.Role, error) {
	//TODO implement me
	panic("implement me")
}

func NewRoleRepository(db *mongo.Database, collectionName string) role.RoleRepository {
	baseRepo := contracts.NewBaseRepository[domain.Role](db, collectionName).(*contracts.BaseRepository[domain.Role])
	return &roleRepository{BaseRepository: *baseRepo }
}


//func (rr *roleRepository) GetRole(ctx context.Context, id string) (*domain.Role, error) {
//	collection := rr.db.Database(rr.cfg.DatabaseName).Collection(rr.collectionName)
//	var role domain.Role
//	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role)
//	if err != nil {
//		if errors.Is(err, mongo.ErrNoDocuments) {
//			return nil, nil
//		}
//		return nil, err
//	}
//	return &role, nil
//}
//func (rr *roleRepository) CreateRole(ctx context.Context, role *domain.Role) error {
//	collection := rr.db.Database(rr.cfg.DatabaseName).Collection(rr.collectionName)
//	_, err := collection.InsertOne(ctx, role)
//	return err
//}
//func (rr *roleRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
//	collection := rr.db.Database(rr.cfg.DatabaseName).Collection(rr.collectionName)
//	filter := bson.M{"_id": role.Id}
//	_, err := collection.ReplaceOne(ctx,filter, role)
//	return err
//}
//
//func (rr *roleRepository) DeleteRole(ctx context.Context, id string, deleterId string) error {
//	collection := rr.db.Database(rr.cfg.DatabaseName).Collection(rr.collectionName)
//	filter := bson.M{"_id": id}
//	update := bson.M{"$set":
//		bson.M{ "DeletionTime": time.Now(), "IsDeleted": true, "DeleterId": deleterId },}
//
//	_, err := collection.UpdateOne(ctx,filter,update)
//	return err
//}
//
//func (rr *roleRepository) GetAllRole(ctx context.Context, query contracts.Query) ([]domain.Role, error) {
//	collection := rr.db.Database(rr.cfg.DatabaseName).Collection(rr.collectionName)
//
//	filter, options := query.QueryBuilder()
//	cursor, err := collection.Find(ctx, filter, options)
//
//	if err != nil {
//		return nil, err
//	}
//	defer cursor.Close(ctx)
//
//	var roles []domain.Role
//	if err := cursor.All(ctx, &roles); err != nil {
//		return nil, err
//	}
//
//	return roles, nil
//}


//func NewRoleRepository(db *mongo.Client, cfg *config.Configuration) role.RoleRepository {
//	return &roleRepository{db: db, cfg: cfg, collectionName: "Roles" }
//}