package envvariable

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type GodotEnvVariableService struct{}

func (g GodotEnvVariableService) Load(filenames ...string) error {
	return godotenv.Load(filenames...)
}

func (g GodotEnvVariableService) Get(key string) string {
	result := os.Getenv(key)
	if result == "" {
		log.Printf("Warning: Environment variable '%s' is empty or not set", key)
	}

	return result
}
