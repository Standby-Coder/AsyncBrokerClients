package models

import (
	"fmt"
)

var (
	Config      = &BrokerConfig{}
	AppStateVar = &AppState{}
)

type AppState struct {
	Debug          bool
	ContextTimeout int
}

// For Kafka, we use TopicName instead of QueueName
// Moreover, we don't need Username, Password, or PasswordFile for basic setups
// If there is Kerberos or SASL with TLS, those would be added here
type BrokerConfig struct {
	Host      string
	Port      string
	TopicName string
}

func (bc *BrokerConfig) GetConnectionString() string {
	// For Kafka we typically use bootstrap servers in the form host:port
	if bc.Host == "" || bc.Port == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", bc.Host, bc.Port)
}
