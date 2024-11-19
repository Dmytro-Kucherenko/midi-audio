package config

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Schema struct {
	BindsPath   string `validate:"required,endswith=.json"`
	HistoryPath string `validate:"required,endswith=.json"`
}

var schema Schema

func init() {
	godotenv.Load(".env")
	validate := validator.New()

	schema = Schema{
		BindsPath:   getEnv("BINDS_PATH"),
		HistoryPath: getEnv("HISTORY_PATH"),
	}

	if err := validate.Struct(schema); err != nil {
		panic(fmt.Sprintln("Enviromental variables error:", err.Error()))
	}
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func BindsPath() string {
	return schema.BindsPath
}

func HistoryPath() string {
	return schema.HistoryPath
}
