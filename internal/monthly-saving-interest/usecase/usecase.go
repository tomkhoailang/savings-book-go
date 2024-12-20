package usecase

import (
	"context"
	"errors"

	"SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/domain"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	saving_book "SavingBooks/internal/saving-book"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type monthlyUC struct {
	monthlyRepo  monthly_saving_interest.Repository
	savingRepo saving_book.SavingBookRepository

}

func (m *monthlyUC) GetTotalEarningsOfSavingBooks(ctx context.Context, savingBookIds []string) (map[string]float64, error) {
	collectionInterface := m.monthlyRepo.GetCollection()
	collection := collectionInterface.(*mongo.Collection)

	objectIDs := make([]primitive.ObjectID, 0, len(savingBookIds))
	for _, id := range savingBookIds {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(err.Error())
		}
		objectIDs = append(objectIDs, objectID)
	}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"SavingBookId": bson.M{
					"$in": objectIDs,
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "SavingBook",
				"localField":   "SavingBookId",
				"foreignField": "_id",
				"as":           "saving_book",
			},
		},
		{
			"$group": bson.M{
				"_id": "$SavingBookId",
				"totalAmount": bson.M{"$sum": "$Amount"},
			},
		},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make(map[string]float64)
	for cursor.Next(ctx) {
		var result struct {
			SavingBookId primitive.ObjectID `bson:"_id"`
			TotalAmount  float64            `bson:"totalAmount"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results[result.SavingBookId.Hex()] = result.TotalAmount
	}

	return results, nil
}

func (m *monthlyUC) GetListMonthlyInterest(ctx context.Context, query *contracts.Query, auth *presenter.AuthData) (*contracts.QueryResult[domain.MonthlySavingInterest], error) {
	var interestInterface interface{}
	var err error

	if _, ok := auth.Roles["Admin"]; ok {
		interestInterface, err = m.monthlyRepo.GetList(ctx, query)
	} else {
		interestInterface, err = m.monthlyRepo.GetListAuth(ctx, query, auth.UserId)
	}

	if err != nil {
		return nil, err
	}

	interest := interestInterface.(*contracts.QueryResult[domain.MonthlySavingInterest])
	return interest, nil
}

func (m *monthlyUC) GetListMonthlyInterestOfSavingBook(ctx context.Context, query *contracts.Query, userId, savingBookId string) (*contracts.QueryResult[domain.MonthlySavingInterest], error) {
	var monthlyInterfaces interface{}
	var err error
	monthlyInterfaces, err = m.monthlyRepo.GetListAuthOnReference(ctx, query, "", "SavingBookId", savingBookId)
	if err != nil {
		return nil, err
	}

	monthly := monthlyInterfaces.(*contracts.QueryResult[domain.MonthlySavingInterest])
	return monthly, nil
}

func NewMonthlyUC(monthlyRepo monthly_saving_interest.Repository, savingRepo saving_book.SavingBookRepository) monthly_saving_interest.UseCase {
	return &monthlyUC{monthlyRepo: monthlyRepo, savingRepo: savingRepo}
}
