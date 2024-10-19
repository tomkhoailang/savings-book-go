package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	presenter3 "SavingBooks/internal/auth/presenter"
	"SavingBooks/internal/contracts"
	"SavingBooks/internal/contracts/paypal"
	"SavingBooks/internal/domain"
	"SavingBooks/internal/notification"
	presenter2 "SavingBooks/internal/notification/presenter"
	"SavingBooks/internal/payment"
	saving_book "SavingBooks/internal/saving-book"
	"SavingBooks/internal/saving-book/presenter"
	regulation "SavingBooks/internal/saving-regulation"
	kafka2 "SavingBooks/internal/services/kafka"
	"SavingBooks/internal/services/kafka/event"
	transaction_ticket "SavingBooks/internal/transaction-ticket"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type savingBookUseCase struct {
	savingBookRepo saving_book.SavingBookRepository
	regulationRepo regulation.SavingRegulationRepository
	ticketRepo     transaction_ticket.TransactionTicketRepository
	paymentUC      payment.PaymentUseCase
	notificationUC notification.UseCase
	kafkaProducer  *kafka2.KafkaProducer
}

func (s *savingBookUseCase) WithdrawOnline(ctx context.Context, input *presenter.WithDrawInput, savingBookId, userId string) error {
	savingBook, err := s.savingBookRepo.Get(ctx, savingBookId)
	if err != nil {
		return err
	}
	if savingBook.CreatorId.Hex() != userId {
		return errors.New(saving_book.NotSavingBookOwnerError)
	}
	if savingBook.Balance < input.Amount {
		return errors.New(saving_book.InsufficientBalance)
	} else if savingBook.Regulations[len(savingBook.Regulations)-1].MinWithDrawValue > input.Amount {
		return errors.New(saving_book.MinWithdrawValueError)
	}
	savingBook.Balance -= input.Amount

	ticket := &domain.TransactionTicket{
		SavingBookId:    savingBook.Id,
		TransactionDate: time.Now(),
		Status:          transaction_ticket.TransactionStatusPending,
		Email:           input.Email,
		PaymentType:     saving_book.TransactionTypeWithdraw,
		PaymentAmount:   input.Amount,
	}
	ticket.SetCreate(userId)

	session, err := s.ticketRepo.GetMongoClient().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		_, err = s.savingBookRepo.Update(sessionContext, savingBook,savingBook.Id.Hex(), []string{"Balance"})
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return err
		}
		err = s.ticketRepo.Create(sessionContext, ticket)
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return err
		}
		return session.CommitTransaction(sessionContext)
	})

	if err != nil {
		return err
	}

	withDrawMessage := &event.WithDrawEvent{
		Amount:        input.Amount,
		SavingBookId:  savingBookId,
		TransactionId: ticket.Id.Hex(),
		Email:         ticket.Email,
	}

	eventJson, err := json.Marshal(withDrawMessage)
	err = s.kafkaProducer.SendMessage(kafka2.CaptureOrderTopic, eventJson)
	if err != nil {
		return err
	}
	return nil
}

func (s *savingBookUseCase) HandleWithdraw(ctx context.Context, input *event.WithDrawEvent) error {
	ticket, err := s.ticketRepo.Get(ctx, input.TransactionId)
	if err != nil {
		return err
	}
	ticket.Status = transaction_ticket.TransactionStatusSuccess

	payoutRequest := &paypal.UCPayoutRequest{
		Amount: input.Amount,
		Email:  input.Email,
	}

	session, err := s.ticketRepo.GetMongoClient().StartSession()
	if err != nil {
		return nil
	}
	defer session.EndSession(ctx)


	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		_, err = s.paymentUC.SendPayout(sessionContext, payoutRequest)
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return err
		}

		_, err = s.ticketRepo.Update(sessionContext, ticket, ticket.Id.Hex(), []string{"Status"})
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return err
		}
		return nil
	})
	if err != nil {
		if revertErr := s.revertBalanceAndNotify(ctx, input, ticket); revertErr != nil {
			log.Printf("Failed to revert balance and notify: %v", err)
		}
		return err
	}

	notification := &presenter2.NotificationInput{
		UserId:              ticket.CreatorId,
		Message:             "Withdraw successfully",
		TransactionTicketId: ticket.Id,
		Status:              transaction_ticket.TransactionStatusSuccess,
	}
	err = s.notificationUC.SendNotification(ctx, notification)
	if err != nil {
		log.Printf("Failed to send success notification: %v", err)
		return err
	}
	return nil
}

func (s *savingBookUseCase) ConfirmPaymentOnline(ctx context.Context, paymentId ,userId string) error {

	ticket, err := s.ticketRepo.GetByField(ctx, "PaymentId", paymentId)
	if err != nil {
		return err
	}
	if ticket.CreatorId.Hex() != userId {
		return errors.New(saving_book.NotSavingBookOwnerError)
	}

	savingBook, err := s.savingBookRepo.Get(ctx, ticket.SavingBookId.Hex())
	if err != nil {
		return err
	}

	if ticket.Status != transaction_ticket.TransactionStatusPending {
		return errors.New(saving_book.TransactionTicketNotPendingStatus)
	}

	resp, err := s.paymentUC.CaptureOrder(ctx, paymentId)
	if err != nil {
		return err
	}
	session, err := s.ticketRepo.GetMongoClient().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		ticket.Status = transaction_ticket.TransactionStatusSuccess
		ticket.PaymentLink = ""
		ticket.Email = resp.Payer.EmailAddress
		ticket.TransactionDate = time.Now()
		ticket.SetSysUpdate()


		nextMonth  := time.Now().AddDate(0, 1, 1)
		savingBook.Balance += ticket.PaymentAmount
		updateFields := []string{"Balance", "NextScheduleMonth"}

		if savingBook.Status != saving_book.SavingBookActive {
			updateFields = append(updateFields, "Status")
			savingBook.Status = saving_book.SavingBookActive
		}

		savingBook.NextScheduleMonth = time.Date(nextMonth.Year(), nextMonth.Month(), nextMonth.Day(), 0, 0, 0, 0, time.Local)

		_, err = s.savingBookRepo.Update(ctx, savingBook, savingBook.Id.Hex(), updateFields)

		if err != nil {
			_ = session.AbortTransaction(ctx)
			return err
		}
		_, err = s.ticketRepo.Update(ctx, ticket, ticket.Id.Hex(),[]string{"Status", "PaymentLink", "TransactionDate", "Email"})
		if err != nil {
			_ = session.AbortTransaction(ctx)
			return err
		}
		return session.CommitTransaction(ctx)
	})
	return nil
}

func (s *savingBookUseCase) CreateSavingBookOnline(ctx context.Context, input *presenter.SavingBookGuestInput, creatorId string) (*domain.SavingBook, error) {
	entity := &domain.SavingBook{}
	err := copier.Copy(entity, &input)
	if err != nil {
		return nil, err
	}
	regulation, err := s.regulationRepo.Get(ctx, input.RegulationId)
	if err != nil {
		return nil, err
	}

	selectedReq, err := findSavingType(input.Term, regulation)
	if err != nil {
		return nil, err
	}
	noTermReq, err := findSavingType(0, regulation)
	if err != nil {
		return nil, err
	}

	req := &domain.Regulation{
		RegulationIdRef:  regulation.Id,
		ApplyDate:        time.Now(),
		Name:             selectedReq.Name,
		TermInMonth:      selectedReq.Term,
		InterestRate:     selectedReq.InterestRate,
		MinWithDrawValue: regulation.MinWithdrawValue,
		MinWithDrawDay:   regulation.MinWithdrawDay,
		NoTermInterestRate: noTermReq.InterestRate,
	}

	entity.Regulations = append(entity.Regulations, *req)
	entity.Status = saving_book.SavingBookInit

	entity.SetCreate(creatorId)

	objectId, err := primitive.ObjectIDFromHex(creatorId)
	if err != nil {
		return nil, err
	}
	entity.AccountId = objectId

	resp, err := s.paymentUC.CreateOrder(ctx, &paypal.InitOrderRequest{
		SavingBookId: entity.Id.Hex(),
		Amount:       fmt.Sprintf("%.2f", input.NewPaymentAmount),
	})
	if err != nil {
		return nil, err
	}

	session, err := s.ticketRepo.GetMongoClient().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			return err
		}
		err = s.savingBookRepo.Create(ctx, entity)
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return err
		}
		ticket := &domain.TransactionTicket{
			SavingBookId:    entity.Id,
			TransactionDate: time.Time{},
			Status:          transaction_ticket.TransactionStatusPending,
			Email:           "",
			PaymentLink:     resp.Links[1].Href,
			PaymentType:     saving_book.TransactionTypeDeposit,
			PaymentId:       resp.Id,
			PaymentAmount:   input.NewPaymentAmount,
		}
		ticket.SetCreate(creatorId)

		err = s.ticketRepo.Create(sessionContext, ticket)
		if err != nil {
			_ = session.AbortTransaction(sessionContext)
			return err
		}
		return session.CommitTransaction(sessionContext)
	})

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *savingBookUseCase) GetListSavingBook(ctx context.Context, query *contracts.Query, auth *presenter3.AuthData) (*contracts.QueryResult[domain.SavingBook], error) {
	var savingBookInterfaces interface{}
	var err error

	if _, ok := auth.Roles["Admin"]; ok {
		savingBookInterfaces, err = s.savingBookRepo.GetList(ctx, query)
	} else {
		savingBookInterfaces, err = s.savingBookRepo.GetListAuth(ctx, query, auth.UserId)
	}

	if err != nil {
		return nil, err
	}
	savingBooks := savingBookInterfaces.(*contracts.QueryResult[domain.SavingBook])

	return savingBooks, nil
}

func (s *savingBookUseCase) revertBalanceAndNotify(ctx context.Context, input *event.WithDrawEvent, ticket *domain.TransactionTicket) error {
	savingBook, err := s.savingBookRepo.Get(ctx, ticket.SavingBookId.Hex())
	if err != nil {
		return fmt.Errorf("failed to get saving book: %w", err)
	}

	savingBook.Balance += input.Amount
	if _, err := s.savingBookRepo.Update(ctx, savingBook, savingBook.Id.Hex(), []string{"Balance"}); err != nil {
		return fmt.Errorf("failed to revert saving book balance: %w", err)
	}

	notification := &presenter2.NotificationInput{
		UserId:              ticket.CreatorId,
		Message:             "Withdrawal failed. Your balance has been restored.",
		TransactionTicketId: ticket.Id,
		Status:              transaction_ticket.TransactionStatusAbort,
	}
	if err := s.notificationUC.SendNotification(ctx, notification); err != nil {
		return fmt.Errorf("failed to send failure notification: %w", err)
	}

	return nil
}

func NewSavingBookUseCase(regulationRepo regulation.SavingRegulationRepository, savingBookRepo saving_book.SavingBookRepository, ticketRepo transaction_ticket.TransactionTicketRepository, paymentUC payment.PaymentUseCase, notificationUC notification.UseCase, kafkaProducer *kafka2.KafkaProducer) saving_book.UseCase {
	return &savingBookUseCase{regulationRepo: regulationRepo, savingBookRepo: savingBookRepo, ticketRepo: ticketRepo, paymentUC: paymentUC, notificationUC: notificationUC, kafkaProducer: kafkaProducer}
}
