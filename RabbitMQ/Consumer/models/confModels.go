package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Config      = &BrokerConfig{}
	AppStateVar = &AppState{}
)

type AppState struct {
	Debug          bool `mapstructure:"debug"`
	ContextTimeout int  `mapstructure:"contexttimeout"`
}

type BrokerConfig struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Username    string `mapstructure:"username"`
	QueueName   string `mapstructure:"queuename"`
	PasswordEnv string `mapstructure:"passwordenv"`
	Password    string `mapstructure:"-"`
}

func (bc *BrokerConfig) LoadPasswordFromEnv() {
	godotenv.Load("./configs/.env")
	bc.Password = os.Getenv(bc.PasswordEnv)
}

func (bc *BrokerConfig) GetConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", bc.Username, bc.Password, bc.Host, bc.Port)
}
