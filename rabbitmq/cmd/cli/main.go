package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnErr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	// Connect to RabbitMQ and open a channel (same code for input or output)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErr(err, "connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErr(err, "open channel")
	defer ch.Close()

	// Declare a queue (creates it if it doesn't exist)
	_, err = ch.QueueDeclare(
		"my-queue", // name
		true,       // durable
		false,      // autoDelete
		false,      // exclusive
		false,      // noWait
		nil,        // args
	)
	failOnErr(err, "declare queue")

	body := "some text or other, such as JSON"
	err = ch.Publish(
		"",         // exchange (empty = default direct exchange)
		"my-queue", // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnErr(err, "publish message")

	log.Printf("Sent: %s", body)
}
