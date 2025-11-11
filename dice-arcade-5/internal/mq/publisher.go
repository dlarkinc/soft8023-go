package mq

import (
	"log"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// Connect establishes a connection to RabbitMQ.
func Connect(amqpURL string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare an exchange for game events
	err = ch.ExchangeDeclare(
		"game.events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-delete
		false,         // internal
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Publisher{conn: conn, channel: ch}, nil
}

// Publish a simple text message to the exchange with a routing key.
func (p *Publisher) Publish(routingKey, message string) {
	err := p.channel.Publish(
		"game.events", // exchange
		routingKey,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
}

// Close connection when done.
func (p *Publisher) Close() {
	p.channel.Close()
	p.conn.Close()
}
