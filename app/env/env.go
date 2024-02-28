package env

import (
	"fmt"
	"halodeksik-be/app/applogger"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "halodeksik-be"

func LoadEnv() error {
	err := godotenv.Load()
	if err == nil {
		return nil
	}

	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err = godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	return nil
}

func Get(key string) string {
	err := LoadEnv()

	if err != nil {
		applogger.Log.Error(err)
	}

	return os.Getenv(key)
}
