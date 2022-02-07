package env

import (
	"errors"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Env struct {
	Name string
	File string
	Line int
	Url  string
}

func LoadEnvs(config interface{}) error {
	if err := godotenv.Load(); err != nil {
		log.Printf("error while .env reading")
	}
	if err := env.Parse(config); err != nil {
		return err
	}
	return nil
}

func LoadEnvFileIfExists(envFile string) error {
	if envFile == "" {
		envFile = ".env"
	}
	if _, err := os.Stat(envFile); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
	} else {
		if err := godotenv.Load(envFile); err != nil {
			return err
		}
	}
	return nil
}
