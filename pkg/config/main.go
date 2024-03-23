package config

import (
	"log"

	"github.com/Bismyth/game-server/pkg/db"
	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Application struct {
		Production  bool
		BindAddress string `validate:"required"`
	}

	Redis *db.Config `validate:"required"`
}

var (
	k      = koanf.New(".")
	parser = yaml.Parser()
)

func New() *Config {
	var C Config

	if err := k.Load(file.Provider("config.yml"), parser); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if err := k.Unmarshal("", &C); err != nil {
		log.Fatalf("error unmarshalling config: %v", err)
	}

	validate := validator.New()
	err := validate.Struct(C)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Fatalf("unable to read from updated config: %v", validationErrors)
	}

	db.SetConfig(C.Redis)

	return &C
}
