package config

import (
	"os"
	"strconv"
)

func GetEnv(envName string, defaultValue string) string {
	var env string
	if env = os.Getenv(envName); env == "" {
		return defaultValue
	}
	return env
}

func GetEnvInt(envName string, defaultValue int) int {
	var env string
	if env = os.Getenv(envName); env == "" {
		return defaultValue
	}
	res, err := strconv.Atoi(env)
	if err != nil {
		panic(err)
	}
	return res
}

