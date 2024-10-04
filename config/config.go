package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Configuration struct {
	Port                  string `env:"PORT" envDefault:"8080"`
	HashSalt              string `env:"HASH_SALT"`
	TokenDuration         int64  `env:"TOKEN_DURATION"`
	RefreshTokenDuration  int64  `env:"REFRESH_TOKEN_DURATION"`
	JwtSecret             string `env:"JWT_SECRET"`
	DatabaseConnectionURL string `env:"CONNECTION_URL"`
	DatabaseName          string `env:"DB_NAME"`
}

func NewConfig() *Configuration {
	err := godotenv.Load()

	if err != nil {
		log.Println("No .env file could be found: ", err)
	}
	var cfg Configuration
	err = env.Parse(&cfg)
	if err != nil {
		log.Println("Parsing cfg file err: ", err)
	}
	return &cfg
}
