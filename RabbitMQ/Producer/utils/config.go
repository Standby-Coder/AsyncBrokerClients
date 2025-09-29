package utils

import (
	"context"
	"log"
	"producer/models"
	"time"

	"github.com/spf13/viper"
)

var (
	Ctx, Cancel = getContext()
)

func InitConfig() {
	// Set up Viper to read the config file
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// This reads the config file and then returns a viper
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	log.Println("Config file loaded successfully")

	// Map the config values to the BrokerConfig struct
	if err := viper.UnmarshalKey("BrokerConfig", models.Config); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}
	models.Config.LoadPasswordFromEnv()

	// Map the config values to the AppState struct
	if err := viper.UnmarshalKey("AppConfig", models.AppStateVar); err != nil {
		log.Fatalf("Error unmarshaling app state config: %v", err)
	}

	Info("Configuration initialized successfully")
	Info("Debug mode is enabled")
	Info("Broker configuration loaded successfully")
	Info("Connection string: " + models.Config.GetConnectionString())
}

func getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(),
		time.Duration(models.AppStateVar.ContextTimeout)*time.Second)
}
