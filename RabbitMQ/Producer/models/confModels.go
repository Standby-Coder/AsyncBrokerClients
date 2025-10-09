package models

import (
	"encoding/base64"
	"fmt"
	"os"
)

var (
	Config      = &BrokerConfig{}
	AppStateVar = &AppState{}
)

type AppState struct {
	Debug          bool
	ContextTimeout int
}

type BrokerConfig struct {
	Host         string
	Port         string
	Username     string
	QueueName    string
	Password     string
	PasswordFile string
}

func (bc *BrokerConfig) LoadPasswordFromFile() {
	data, err := os.ReadFile(bc.PasswordFile)
	if err != nil {
		fmt.Printf("Failed to read password file: %v", err)
	}
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		fmt.Printf("Failed to decode password from file: %v", err)
	}
	bc.Password = string(decoded)
}

func (bc *BrokerConfig) GetConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/", bc.Username, bc.Password, bc.Host, bc.Port)
}
