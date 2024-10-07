package server

import (
	"context"
	"fmt"

	"SavingBooks/internal/auth/middleware"
	"github.com/gin-gonic/gin"

	authHttp "SavingBooks/internal/auth/delivery/http"
	authRepo "SavingBooks/internal/auth/repository"
	authUC "SavingBooks/internal/auth/usecase"

	testHttp "SavingBooks/internal/not-auth/delivery/http"

	roleHttp "SavingBooks/internal/role/delivery/http"
	roleRepo "SavingBooks/internal/role/repository"
	roleUC "SavingBooks/internal/role/usecase"

	paymentHttp "SavingBooks/internal/payment/delivery/http"
	paymentUC "SavingBooks/internal/payment/usecase"
)

func (s *Server) MapHandlers(g *gin.Engine) error {
	db := s.db.Database(s.cfg.DatabaseName)


	userRepo := authRepo.NewUserRepository(db,"Users")
	roleRepo := roleRepo.NewRoleRepository(db, "Roles")
	ctx := context.Background()
	if err := roleRepo.SeedRole(ctx); err != nil {
		fmt.Println("Something wrong with seed roles")
		return err
	}

	authUC := authUC.NewAuthUseCase(userRepo, roleRepo, s.cfg.HashSalt, []byte(s.cfg.JwtSecret), s.cfg.TokenDuration, s.cfg.RefreshTokenDuration)
	roleUc := roleUC.NewRoleUseCase(roleRepo)
	paymentUC := paymentUC.NewPaymentUseCase(s.cfg.ClientId, s.cfg.ClientSecret)


	authHandler := authHttp.NewAuthHandler(authUC)
	roleHandler := roleHttp.NewRoleHandler(roleUc)
	paymentHandler := paymentHttp.NewPaymentHandler(paymentUC)

	v1 := g.Group("/api/v1")

	authGroup := v1.Group("/auth")
	testGroup := v1.Group("/test")
	roleGroup := v1.Group("/role")
	paymentGroup := v1.Group("/payment")

	mw := middleware.NewMiddleWareManager(authUC)

	authHttp.MapAuthRoutes(authGroup, authHandler, mw)
	testHttp.MapAuthRoutes(testGroup, authHandler, mw)
	roleHttp.MapAuthRoutes(roleGroup, roleHandler, mw)
	paymentHttp.MapAuthRoutes(paymentGroup, paymentHandler)
	return nil
}
