package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"producer/models"
	"producer/utils"
	"syscall"
	"time"

	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

func main() {
	// Initialize configs from ./configs/config.yaml and .env
	utils.InitConfig()

	// Build bootstrap servers string from config
	bootstrap := models.Config.GetConnectionString()
	if bootstrap == "" {
		utils.Fatal("Invalid bootstrap servers in config")
		fmt.Println("Invalid bootstrap servers in config")
		return
	}
	utils.Info(fmt.Sprintf("Bootstrap servers: %s", bootstrap))

	if models.Config.TopicName == "" {
		utils.Fatal("No topic provided in config")
		fmt.Println("No topic provided in config")
		return
	}
	utils.Info(fmt.Sprintf("Using topic : %s", models.Config.TopicName))

	// Create Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrap})
	if err != nil {
		utils.Fatal(fmt.Sprintf("Failed to create Kafka producer: %v", err))
		fmt.Println("Failed to create Kafka producer:", err)
		return
	}
	defer p.Close()
	utils.Info(fmt.Sprintf("Kafka producer created bootstrap=%s topic=%s", bootstrap, models.Config.TopicName))

	// Channel for messages and graceful shutdown
	msgChan := make(chan string)
	done := make(chan struct{})

	// Delivery report handler goroutine
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					utils.Error(fmt.Sprintf("Delivery failed: %v", ev.TopicPartition.Error))
				} else {
					utils.Info(fmt.Sprintf("Delivered message to partition %v", ev.TopicPartition))
				}
			}
		}
	}()

	// Goroutine to publish messages
	go func() {
		for {
			select {
			case msg := <-msgChan:
				kt := &kafka.Message{
					TopicPartition: kafka.TopicPartition{Topic: &models.Config.TopicName, Partition: kafka.PartitionAny},
					Value:          []byte(msg),
				}
				if err := p.Produce(kt, nil); err != nil {
					utils.Error(fmt.Sprintf("Failed to produce message: %v", err))
				} else {
					utils.Info(fmt.Sprintf("Message queued for delivery: %s", msg))
				}
			case <-done:
				fmt.Println("Shutting down producer...")
				// Wait for message deliveries
				p.Flush(5000)
				return
			}
		}
	}()

	// Capture interrupt signals to allow graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Read user input in the main goroutine
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages to send to Kafka. Press Ctrl+C to exit.")
	for {
		fmt.Print("> ")
		select {
		case sig := <-sigChan:
			fmt.Println("Exiting application")
			utils.Info(fmt.Sprintf("Received signal: %v", sig))
			close(done)

			// allow goroutine to free up memory
			time.Sleep(200 * time.Millisecond)
			return
		default:
			if !scanner.Scan() {
				close(done)
				return
			}
			input := scanner.Text()
			if input == "exit" {
				utils.Info("Exiting application on user command")
				close(done)

				//Same here as above
				time.Sleep(200 * time.Millisecond)
				return
			}
			// Send input to publishing goroutine
			msgChan <- input
			utils.Info(fmt.Sprintf("Input received: %s", input))
		}
	}
}
