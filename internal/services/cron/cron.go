package cron

import (
	"context"
	"errors"
	"log"
	"time"

	"SavingBooks/internal/domain"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	saving_book "SavingBooks/internal/saving-book"
	"github.com/robfig/cron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Scheduler struct {
	savingBookRepo          saving_book.SavingBookRepository
	monthSavingInterestRepo monthly_saving_interest.Repository
	cron                    *cron.Cron
}

func NewScheduler(sv saving_book.SavingBookRepository, monthSavingInterestRepo monthly_saving_interest.Repository) *Scheduler {
	return &Scheduler{savingBookRepo: sv, monthSavingInterestRepo: monthSavingInterestRepo, cron: cron.New(cron.WithSeconds())}
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
	savingBookInterface := s.savingBookRepo.GetCollection()
	savingBookCollection := savingBookInterface.(*mongo.Collection)

	now := time.Now()
	filterDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	filter := bson.M{
		"NextScheduleMonth": filterDate,
		"Balance":           bson.M{"$gt": 0},
	}

	cursor, err := savingBookCollection.Find(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return
	}
	defer cursor.Close(context.Background())
	var savingOperations []mongo.WriteModel
	var monthlyOperations []mongo.WriteModel
	batchSize := 50

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
		newBalance := savingBook.Balance * (1 + (interestRate / 100))
		updateDoc["Balance"] = newBalance

		savingBookUpdate := mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": savingBook.Id}).
			SetUpdate(bson.M{"$set": updateDoc})
		savingOperations = append(savingOperations, savingBookUpdate)

		monthlyInterest := domain.MonthlySavingInterest{
			SavingBookId: savingBook.Id,
			Amount:       newBalance - savingBook.Balance,
			InterestRate: interestRate,
		}
		monthlyInterest.SetInit()
		monthlyUpdate := mongo.NewInsertOneModel().SetDocument(&monthlyInterest)
		monthlyOperations = append(monthlyOperations, monthlyUpdate)

		if len(savingOperations) == batchSize {
			err := bulkWrite(s.monthSavingInterestRepo, s.savingBookRepo, savingOperations, monthlyOperations)
			if err != nil {
				log.Println("error after maximum retry:", err)
				return
			}
			savingOperations = savingOperations[:0]
		}
	}
	if len(savingOperations) > 0 {
		err := bulkWrite(s.monthSavingInterestRepo, s.savingBookRepo, savingOperations, monthlyOperations)
		if err != nil {
			log.Println("error after maximum retry:", err)
			return
		}
	}

	if err = cursor.Err(); err != nil {
		log.Println("Error iterating cursor:", err)
	}

}
func bulkWrite(monthlyRepo monthly_saving_interest.Repository, savingBookRepo saving_book.SavingBookRepository, writeSavingBook, writeMonthly []mongo.WriteModel) error {

	return retry(5, 10*time.Second, func() error {
		savingBookInterface := savingBookRepo.GetCollection()
		savingBookCollection := savingBookInterface.(*mongo.Collection)

		monthlyInterface := monthlyRepo.GetCollection()
		monthlyCollection := monthlyInterface.(*mongo.Collection)

		session, err := monthlyRepo.GetMongoClient().StartSession()
		if err != nil {
			log.Println("error starting mongo session:", err)
			return errors.New("error starting mongo session")
		}

		defer session.EndSession(context.Background())
		err = mongo.WithSession(context.Background(), session, func(sessionContext mongo.SessionContext) error {
			if err = session.StartTransaction(); err != nil {
				return err
			}
			_, err = savingBookCollection.BulkWrite(sessionContext, writeSavingBook)
			if err != nil {
				log.Println("Error in BulkWrite:", err)
			}
			if err != nil {
				_ = session.AbortTransaction(sessionContext)
				return err
			}
			_, err = monthlyCollection.BulkWrite(sessionContext, writeMonthly)
			if err != nil {
				log.Println("Error in BulkWrite:", err)
			}
			if err != nil {
				_ = session.AbortTransaction(sessionContext)
				return err
			}
			return session.CommitTransaction(sessionContext)
		})
		if err != nil {
			log.Println("bulk write from cron failed", err)
			return errors.New("bulk write from cron failed")
		}
		return nil
	})

}
func retry(attempts int, sleep time.Duration, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(sleep * time.Duration(i+1))
	}
	return err
}
func monthsBetween(start, end time.Time) int {
	yearDiff := end.Year() - start.Year()
	monthDiff := int(end.Month()) - int(start.Month())

	return yearDiff*12 + monthDiff
}
