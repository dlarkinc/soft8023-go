package main

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Open a channel
	ch, err := conn.Channel()
	failOnErr(err, "Failed to open a channel")
	defer ch.Close()

	// Declare a fanout exchange
	err = ch.ExchangeDeclare(
		"logs",   // exchange name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnErr(err, "Failed to declare exchange")

	// Publish a test message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "info: Hello World!"

	fmt.Printf("About to publish message '%s'\n", body)

	err = ch.PublishWithContext(ctx,
		"logs", // exchange name
		"",     // routing key ignored by fanout
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnErr(err, "Failed to publish message")

	log.Printf("[x] Sent: %s", body)
}
