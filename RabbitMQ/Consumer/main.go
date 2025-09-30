package main

import (
	"bufio"
	"consumer/models"
	"consumer/utils"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Initialize configs from ./configs/config.yaml and .env
	utils.InitConfig()

	// Establish connection to RabbitMQ
	conn, err := amqp.Dial(models.Config.GetConnectionString())
	if err != nil {
		utils.Fatal(fmt.Sprintf("Failed to connect to RabbitMQ: %v", err))
		fmt.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()
	utils.Info("Successfully connected to RabbitMQ")

	// Open a channel
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Failed to open a channel:", err)
		utils.Fatal(fmt.Sprintf("Failed to open a channel: %v", err))
		return
	}
	defer ch.Close()

	// Declare a queue, this creates the queue if it doesn't exist (idempotent function)
	_, err = ch.QueueDeclare(
		models.Config.QueueName,
		true,  // durable
		false, //autoDelete
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Failed to declare a queue: %v", err))
		fmt.Println("Failed to declare a queue:", err)
		return
	}
	utils.Info(fmt.Sprintf("Queue declared: %s", models.Config.QueueName))

	// Creating interactive session for sending messages

	// Set up consumer
	msgs, err := ch.Consume(
		models.Config.QueueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		utils.Fatal(fmt.Sprintf("Failed to register a consumer: %v", err))
		fmt.Println("Failed to register a consumer:", err)
		return
	}

	// Channel for graceful shutdown
	done := make(chan struct{})

	// Goroutine to handle incoming messages
	go func() {
		for {
			select {
			case d, ok := <-msgs:
				if !ok {
					fmt.Println("Message channel closed")
					return
				}
				fmt.Printf("Received message: %s\n", d.Body)
				utils.Info(fmt.Sprintf("Message consumed: %s", d.Body))
			case <-done:
				fmt.Println("Shutting down consumer...")
				return
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
			return
		default:
			if scanner.Scan() {
				input := scanner.Text()
				if input == "exit" {
					utils.Info("Exiting application on user command")
					close(done)
					return
				}
			}
		}
	}
}
