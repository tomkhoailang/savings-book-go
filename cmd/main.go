package main

import (
	"context"
	"log"
	"time"

	"SavingBooks/config"
	"SavingBooks/internal/server"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
// @title Saving Books API
func main() {
	cfg := config.NewConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DatabaseConnectionURL))
	if err != nil {
		log.Fatal("Error mongo connection:", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}
	log.Println("Successfully connect to", cfg.DatabaseConnectionURL)

	s := server.NewServer(cfg, client, logrus.New(),nil)

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}