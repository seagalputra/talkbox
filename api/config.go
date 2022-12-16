package api

import (
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

type (
	appConfig struct {
		DatabaseURL          string `env:"DATABASE_URL"`
		DatabaseName         string `env:"DATABASE_NAME"`
		JwtSecret            string `env:"JWT_SECRET"`
		ServerPort           string `env:"SERVER_PORT"`
		RedisHost            string `env:"REDIS_HOST"`
		SMTPHost             string `env:"SMTP_HOST"`
		SMTPPort             string `env:"SMTP_PORT"`
		SMTPUsername         string `env:"SMTP_USERNAME"`
		SMTPPassword         string `env:"SMTP_PASSWORD"`
		EmailSenderName      string `env:"EMAIL_SENDER_NAME"`
		EmailConfirmationURL string `env:"EMAIL_CONFIRMATION_URL"`
		CloudinaryCloudName  string `env:"CLOUDINARY_CLOUD_NAME"`
		CloudinaryAPIKey     string `env:"CLOUDINARY_API_KEY"`
		CloudinaryAPISecret  string `env:"CLOUDINARY_API_SECRET"`
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
