package server

import (
	"SavingBooks/internal/auth/middleware"
	"github.com/gin-gonic/gin"

	authHttp "SavingBooks/internal/auth/delivery/http"
	authRepo "SavingBooks/internal/auth/repository"
	authUC "SavingBooks/internal/auth/usecase"

	testHttp "SavingBooks/internal/not-auth/delivery/http"


	roleHttp "SavingBooks/internal/role/delivery/http"
	roleRepo "SavingBooks/internal/role/repository"
	roleUC "SavingBooks/internal/role/usecase"
)

func (s *Server) MapHandlers(g *gin.Engine) error {
	userRepo := authRepo.NewUserRepository(s.db, s.cfg)
	roleRepo := roleRepo.NewRoleRepository(s.db.Database(s.cfg.DatabaseName), "Roles")

	authUC := authUC.NewAuthUseCase(userRepo, s.cfg.HashSalt, []byte(s.cfg.JwtSecret), s.cfg.TokenDuration)
	roleUc := roleUC.NewRoleUseCase(roleRepo)

	authHandler := authHttp.NewAuthHandler(authUC)
	roleHandler := roleHttp.NewRoleHandler(roleUc)

	v1 := g.Group("/api/v1")

	authGroup := v1.Group("/auth")
	testGroup := v1.Group("/test")
	roleGorup := v1.Group("/role")

	mw := middleware.NewMiddleWareManager(authUC)

	authHttp.MapAuthRoutes(authGroup, authHandler)
	testHttp.MapAuthRoutes(testGroup, authHandler, mw)
	roleHttp.MapAuthRoutes(roleGorup, roleHandler, mw)
	return nil
}
