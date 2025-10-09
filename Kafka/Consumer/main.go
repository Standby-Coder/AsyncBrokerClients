package main

import (
	"bufio"
	"consumer/models"
	"consumer/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	// Initialize configs from ./configs/config.yaml and .env
	utils.InitConfig()

	// Build bootstrap servers
	bootstrap := models.Config.GetConnectionString()
	if bootstrap == "" {
		utils.Fatal("No bootstrap servers configured")
		fmt.Println("No bootstrap servers configured")
		return
	}

	// Topic
	topic := models.Config.TopicName
	if topic == "" {
		utils.Fatal("No topic configured")
		fmt.Println("No topic configured")
		return
	}

	// Create Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrap,
		"group.id":          "async-broker-consumer",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		utils.Fatal(fmt.Sprintf("Failed to create Kafka consumer: %v", err))
		fmt.Println("Failed to create Kafka consumer:", err)
		return
	}
	defer c.Close()
	utils.Info(fmt.Sprintf("Connected to Kafka bootstrap=%s topic=%s", bootstrap, topic))

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		utils.Fatal(fmt.Sprintf("Failed to subscribe to topic %s: %v", topic, err))
		fmt.Println("Failed to subscribe to topic:", err)
		return
	}

	// Channel for graceful shutdown
	done := make(chan struct{})

	// Goroutine to handle incoming messages
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Shutting down consumer goroutine...")
				return
			default:
				ev := c.Poll(500)
				if ev == nil {
					continue
				}
				switch e := ev.(type) {
				case *kafka.Message:
					fmt.Printf("Received message: %s\n", string(e.Value))
					utils.Info(fmt.Sprintf("Message consumed: %s", string(e.Value)))
				case kafka.Error:
					fmt.Printf("Consumer error: %v\n", e)
					utils.Error(fmt.Sprintf("Kafka error: %v", e))
					if e.IsFatal() {
						fmt.Println("Fatal Kafka error, exiting consumer goroutine...")
						return
					}
				}
			}
		}
	}()

	// Capture interrupt signals to allow graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Consuming messages. Press Ctrl+C or type 'exit' to quit.")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		select {
		case sig := <-sigChan:
			fmt.Println("Exiting application")
			utils.Info(fmt.Sprintf("Received signal: %v", sig))
			close(done)

			// give poll loop time to finish
			time.Sleep(300 * time.Millisecond)
			return
		default:
			if scanner.Scan() {
				input := scanner.Text()
				if input == "exit" {
					utils.Info("Exiting application on user command")
					close(done)

					// Same as above
					time.Sleep(300 * time.Millisecond)
					return
				}
			}
		}
	}
}
