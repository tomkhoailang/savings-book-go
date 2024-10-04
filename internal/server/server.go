package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"SavingBooks/config"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	gin *gin.Engine
	cfg *config.Configuration
	db *mongo.Client
	logger *logrus.Logger
	ready chan bool
}

func NewServer( cfg *config.Configuration, db *mongo.Client, logger *logrus.Logger, ready chan bool) *Server {
	return &Server{gin: gin.New(), cfg: cfg, db: db, logger: logger, ready: ready}
}
func (s *Server) Run() error {

	s.gin.Use(gin.Logger())
	s.gin.Use(gin.Recovery())

	server := &http.Server{
		Addr: ":" + s.cfg.Port,
		Handler: s.gin,
		WriteTimeout: time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	if err:= s.MapHandlers(s.gin); err != nil {
		return  err
	}

	if s.ready != nil {
		s.ready <- true
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		//if err := server.ListenAndServeTLS("./self_cert/server.pem", "./self_cert/server.key"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		//	fmt.Printf("Error closed: %s\n", err)
		//}
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error closed: %s\n", err)
		}
	}()
	fmt.Printf("Server is listening on %s\n", server.Addr)
	<-quit
	_, shutdown := context.WithTimeout(context.Background(), time.Second * 2)
	defer shutdown()
	fmt.Println("Server Exited Properly")
	return nil
}
