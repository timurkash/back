package env

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
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
