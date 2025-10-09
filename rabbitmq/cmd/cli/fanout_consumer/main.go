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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErr(err, "Failed to open a channel")
	defer ch.Close()

	// Ensure the fanout exchange exists
	err = ch.ExchangeDeclare(
		"logs",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnErr(err, "Failed to declare exchange")

	// Create a temporary, auto-deleted queue
	q, err := ch.QueueDeclare(
		"",    // empty name = unique random name
		false, // durable
		true,  // auto-delete when consumer disconnects
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnErr(err, "Failed to declare queue")

	// Bind the queue to the exchange (routing key ignored for fanout)
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key ignored
		"logs", // exchange
		false,  // no-wait
		nil,    // args
	)
	failOnErr(err, "Failed to bind queue")

	// Consume messages from this queue
	msgs, err := ch.Consume(
		q.Name,
		"",    // consumer tag
		true,  // auto-ack
		true,  // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnErr(err, "Failed to register consumer")

	log.Println("[*] Waiting for log messages. To exit press CTRL+C")

	for msg := range msgs {
		log.Printf("[x] %s", msg.Body)
	}
}
