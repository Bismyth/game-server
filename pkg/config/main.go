package config

import (
	"log"

	"github.com/Bismyth/game-server/pkg/api"
	"github.com/Bismyth/game-server/pkg/db"
	"github.com/caarlos0/env/v10"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Application struct {
		Production  bool   `env:"PRODUCTION"`
		PublicDir   string `env:"PUBLIC_DIR"`
		BindAddress string `env:"ADDRESS"`
		JWTSecret   string `env:"JWT_SECRET"`
	} `envPrefix:"APPLICATION_"`

	Redis db.Config `validate:"required" envPrefix:"REDIS_"`
}

func New() *Config {
	var C Config

	// Docker production defaults
	C.Application.Production = true
	C.Application.BindAddress = "0.0.0.0:8080"
	C.Application.PublicDir = "/public"

	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file loaded")
	}

	if err := env.Parse(&C); err != nil {
		log.Fatalf("failed to read in config: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(C); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Fatalf("unable to read from updated config: %v", validationErrors)
	}

	api.SetSigningKey(C.Application.JWTSecret)

	db.SetConfig(C.Redis)

	return &C
}
