package server

import (
	"context"
	"fmt"

	"SavingBooks/internal/auth/middleware"
	monthly_saving_interest "SavingBooks/internal/monthly-saving-interest"
	saving_book "SavingBooks/internal/saving-book"
	kafka2 "SavingBooks/internal/services/kafka"
	"SavingBooks/internal/services/websocket"
	"github.com/gin-gonic/gin"

	authHttp "SavingBooks/internal/auth/delivery/http"
	authRepo "SavingBooks/internal/auth/repository"
	authUC "SavingBooks/internal/auth/usecase"

	testHttp "SavingBooks/internal/test-service/delivery/http"
	testUC "SavingBooks/internal/test-service/usecase"

	roleHttp "SavingBooks/internal/role/delivery/http"
	roleRepo "SavingBooks/internal/role/repository"
	roleUC "SavingBooks/internal/role/usecase"

	paymentHttp "SavingBooks/internal/payment/delivery/http"
	paymentUC "SavingBooks/internal/payment/usecase"

	regulationHttp "SavingBooks/internal/saving-regulation/delivery/http"
	regulationRepo "SavingBooks/internal/saving-regulation/repository"
	regulationUC "SavingBooks/internal/saving-regulation/usecase"

	savingBookHttp "SavingBooks/internal/saving-book/delivery/http"
	savingBookRepo "SavingBooks/internal/saving-book/repository"
	savingBookUC "SavingBooks/internal/saving-book/usecase"

	ticketHttp "SavingBooks/internal/transaction-ticket/delivery/http"
	ticketRepo "SavingBooks/internal/transaction-ticket/repository"
	ticketUc "SavingBooks/internal/transaction-ticket/usecase"

	notificationRepo "SavingBooks/internal/notification/repository"
	notificationUC "SavingBooks/internal/notification/usecase"

	monthlyRepo "SavingBooks/internal/monthly-saving-interest/repository"
	monthlyUC "SavingBooks/internal/monthly-saving-interest/usecase"

)

func (s *Server) MapHandlers(g *gin.Engine) (saving_book.UseCase, saving_book.SavingBookRepository, monthly_saving_interest.Repository, error) {
	db := s.db.Database(s.cfg.DatabaseName)


	userRepo := authRepo.NewUserRepository(db,"Users")
	monthlyRepo := monthlyRepo.NewMonthlySavingInterestRepository(s.db,db,"MonthlySavingInterest",)
	roleRepo := roleRepo.NewRoleRepository(db, "Roles")
	regulationRepo := regulationRepo.NewSavingRepository(db, "Regulations")
	savingBookRepo := savingBookRepo.NewSavingBookRepository(db, "SavingBook")
	ticketRepo := ticketRepo.NewTransactionTicketRepository(s.db,db, "TransactionTickets")
	notificationRepo := notificationRepo.NewNotificationRepository(db, "Notifications")



	kafkaProducer := kafka2.NewKafkaProducer("localhost:9092")


	testUC := testUC.NewTestServiceUseCase(kafkaProducer, s.hub)
	authUC := authUC.NewAuthUseCase(userRepo, roleRepo, s.cfg.HashSalt, []byte(s.cfg.JwtSecret), s.cfg.TokenDuration, s.cfg.RefreshTokenDuration)
	roleUc := roleUC.NewRoleUseCase(roleRepo)
	paymentUC := paymentUC.NewPaymentUseCase(s.cfg.ClientId, s.cfg.ClientSecret)
	regulationUC := regulationUC.NewSavingRegulationUseCase(regulationRepo)
	notificationUC := notificationUC.NewNotificationUseCase(notificationRepo, s.hub)
	savingBookUC := savingBookUC.NewSavingBookUseCase(regulationRepo,savingBookRepo,ticketRepo, paymentUC,notificationUC, kafkaProducer)
	ticketUC := ticketUc.NewTransactionTicketUseCase(ticketRepo, savingBookRepo)
	monthlyUC := monthlyUC.NewMonthlyUC(monthlyRepo)


	testHandler := testHttp.NewTestServiceHandler(testUC)
	authHandler := authHttp.NewAuthHandler(authUC)
	roleHandler := roleHttp.NewRoleHandler(roleUc)
	paymentHandler := paymentHttp.NewPaymentHandler(paymentUC)
	regulationHandler := regulationHttp.NewSavingRegulationHandler(regulationUC)
	savingBookHandler := savingBookHttp.NewSavingBookHandler(savingBookUC, ticketUC, monthlyUC)
	ticketHandler := ticketHttp.NewTransactionTicketHandler(ticketUC)




	v1 := g.Group("/api/v1")
	authGroup := v1.Group("/auth")
	testGroup := v1.Group("/test")
	roleGroup := v1.Group("/role")
	paymentGroup := v1.Group("/payment")
	regulationGroup := v1.Group("/regulation")
	savingBookGroup := v1.Group("/saving-book")
	ticketGroup := v1.Group("/transaction-ticket")

	socketGroup := v1.Group("/ws")


	mw := middleware.NewMiddleWareManager(authUC)


	authHttp.MapAuthRoutes(authGroup, authHandler, mw)
	testHttp.MapAuthRoutes(testGroup, testHandler, mw)
	roleHttp.MapAuthRoutes(roleGroup, roleHandler, mw)
	paymentHttp.MapAuthRoutes(paymentGroup, paymentHandler)
	regulationHttp.MapAuthRoutes(regulationGroup, regulationHandler, mw)
	savingBookHttp.MapAuthRoutes(savingBookGroup, savingBookHandler, mw)
	ticketHttp.MapAuthRoutes(ticketGroup, ticketHandler, mw)

	websocket.MapAuthRoutes(socketGroup, s.hub, mw)


	ctx := context.Background()
	if err := roleRepo.SeedRole(ctx); err != nil {
		fmt.Println("Something wrong with seed roles")
		return savingBookUC,savingBookRepo, monthlyRepo,err
	}


	return savingBookUC,savingBookRepo, monthlyRepo, nil
}
