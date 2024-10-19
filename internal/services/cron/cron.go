package cron

import (
	"context"
	"log"
	"time"

	"SavingBooks/internal/domain"
	saving_book "SavingBooks/internal/saving-book"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Scheduler struct {
	savingBookRepo saving_book.SavingBookRepository
	cron           *cron.Cron
}

func NewScheduler(sv saving_book.SavingBookRepository) *Scheduler {
	return &Scheduler{savingBookRepo: sv, cron: cron.New(cron.WithSeconds())}
}
func (s *Scheduler) Start() {
	//_, err := c.AddFunc("@midnight", s.handleSavingBook)
	_, err := s.cron.AddFunc("* * * * * * ", s.handleSavingBook)
	if err != nil {
		log.Println(err)
		return
	}
	s.cron.Start()
}
func (s *Scheduler) Stop() {
	s.cron.Stop()
}
func (s *Scheduler) handleSavingBook() {
	collectionInterface := s.savingBookRepo.GetCollection()
	collection := collectionInterface.(*mongo.Collection)

	now := time.Now()
	filterDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	filter := bson.M{
		"NextScheduleMonth": filterDate,
		"Balance":           bson.M{"$gt": 0},
	}

	cursor, err := collection.Find(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return
	}
	defer cursor.Close(context.Background())
	var operations []mongo.WriteModel
	batchSize := 500

	for cursor.Next(context.Background()) {
		var savingBook domain.SavingBook
		if err := cursor.Decode(&savingBook); err != nil {
			log.Println("Error decoding saving book:", err)
			continue
		}
		if len(savingBook.Regulations) == 0 {
			log.Println("Cannot find a regulation for this saving book")
			continue
		}
		newestRegulation := savingBook.Regulations[len(savingBook.Regulations)-1]

		monthRange := monthsBetween(newestRegulation.ApplyDate, now)
		interestRate := newestRegulation.InterestRate
		updateDoc := bson.M{
			"NextScheduleMonth": now.AddDate(0, 1, 0).Truncate(24 * time.Hour),
		}

		if monthRange >= newestRegulation.TermInMonth {
			if savingBook.Status != saving_book.SavingBookExpired {
				updateDoc["Status"] = saving_book.SavingBookExpired
			}
			interestRate = newestRegulation.NoTermInterestRate

		}
		newBalance := savingBook.Balance + (savingBook.Balance * (interestRate / 100))
		updateDoc["Balance"] = newBalance

		update := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": savingBook.Id}).
			SetUpdate(bson.M{"$set": updateDoc})
		operations = append(operations, update)

		if len(operations) == batchSize {
			_, err = collection.BulkWrite(context.Background(), operations)
			if err != nil {
				log.Println("Error in BulkWrite:", err)
			}
			operations = operations[:0]
		}
	}
	if len(operations) > 0 {
		_, err = collection.BulkWrite(context.Background(), operations)
		if err != nil {
			log.Println("Error in BulkWrite:", err)
		}
	}

	if err = cursor.Err(); err != nil {
		log.Println("Error iterating cursor:", err)
	}

}
func monthsBetween(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())

	return yearDiff*12 + monthDiff
}
