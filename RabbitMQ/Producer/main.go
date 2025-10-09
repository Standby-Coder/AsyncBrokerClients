package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"producer/models"
	"producer/utils"
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

	// Create channels for messages and graceful shutdown
	msgChan := make(chan string)
	done := make(chan struct{})

	// Goroutine to handle message publishing
	go func() {
		for {
			select {
			case msg := <-msgChan: // When message channel receives a message
				// Publish the message to the queue
				err = ch.PublishWithContext(
					utils.Ctx,
					"", // exchange
					models.Config.QueueName,
					false, // mandatory
					false, // immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(msg),
					},
				)
				if err != nil {
					utils.Fatal(fmt.Sprintf("Failed to publish message: %v", err))
					fmt.Println("Failed to publish message")
				} else {
					utils.Info(fmt.Sprintf("Message published: %s", msg))
				}
			case <-done: // When done channel is closed, exit goroutine
				fmt.Println("Shutting down publisher...")
				return
			}
		}
	}()

	// Capture interrupt signals to allow graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Read user input in the main goroutine
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter messages to send to the queue. Press Ctrl+C to exit.")
	for {
		fmt.Print("> ")
		select {
		case sig := <-sigChan: // On interrupt signal from keyboard
			fmt.Println("Exiting application")
			utils.Info(fmt.Sprintf("Received signal: %v", sig))
			close(done)
			return
		default:
			// Read user input
			// If there's an error reading from the scanner, close the done channel
			if !scanner.Scan() {
				close(done)
				return
			}

			input := scanner.Text()
			if input == "exit" { // User can type 'exit' to quit
				utils.Info("Exiting application on user command")
				close(done)
				return
			}

			// Send the input to the message channel for publishing
			msgChan <- input
			utils.Info(fmt.Sprintf("Input received: %s", input))
		}	
	}
}
