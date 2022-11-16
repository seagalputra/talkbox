package config

import (
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type (
	appConfig struct {
		Host               string `env:"HOST"`
		Port               string `env:"PORT"`
		SpreadsheetID      string `env:"SPREADSHEET_ID"`
		ServiceAccountPath string `env:"SERVICE_ACCOUNT_PATH"`
	}
)

var AppConfig appConfig

const tagName = "env"

func LoadAppConfig() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	AppConfig = appConfig{}
	t := reflect.TypeOf(AppConfig)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get(tagName)

		osEnv := os.Getenv(tag)
		reflect.ValueOf(&AppConfig).Elem().FieldByName(field.Name).SetString(osEnv)
	}

	return nil
}
