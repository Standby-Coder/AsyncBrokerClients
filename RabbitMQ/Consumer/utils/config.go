package utils

import (
	"context"
	"os"
	"consumer/models"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	Ctx, Cancel = getContext()
)

func InitConfig() {
	// Get config vars from env file and load into struct
	godotenv.Load(".env")
	viper.AutomaticEnv()

	models.Config = &models.BrokerConfig{
		Host:      viper.GetString("BROKER_HOST"),
		Port:      viper.GetString("BROKER_PORT"),
		Username:  viper.GetString("BROKER_USERNAME"),
		QueueName: viper.GetString("BROKER_QUEUE_NAME"),
	}

	models.AppStateVar = &models.AppState{
		Debug:          viper.GetBool("APP_DEBUG"),
		ContextTimeout: viper.GetInt("APP_CONTEXT_TIMEOUT"),
	}

	// Check if file exists
	if _, err := os.Stat("/run/secrets/broker_password"); err == nil {
		models.Config.PasswordFile = "/run/secrets/broker_password"
	} else {
		models.Config.PasswordFile = "./broker.txt"
	}

	models.Config.LoadPasswordFromFile()

	Info("Configuration initialized successfully")
	Info("Debug mode is enabled")
	Info("Broker configuration loaded successfully")
	Info("Connection string: " + models.Config.GetConnectionString())
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(),
		time.Duration(models.AppStateVar.ContextTimeout)*time.Second)
}
