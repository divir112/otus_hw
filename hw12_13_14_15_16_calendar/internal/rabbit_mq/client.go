package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config IConfig
}

type IConfig interface {
	GetHost() string
	GetPort() string
	GetQueue() string
}

func NewConnection(config IConfig) (*RabbitMQ, error) {
	address := fmt.Sprintf("amqp://guest:guest@%s:%s/", config.GetHost(), config.GetPort())
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, fmt.Errorf("unable to open connect to RabbitMQ server. Error: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("can't open channel. Error: %w", err)
	}

	return &RabbitMQ{
		conn:   conn,
		ch:     ch,
		config: config,
	}, nil
}

func (r *RabbitMQ) Close() error {
	return r.conn.Close()
}

func (r *RabbitMQ) Publish(ctx context.Context, msg []byte) error {
	q, err := r.ch.QueueDeclare(
		r.config.GetQueue(), // name of the queue
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	if err != nil {
		return fmt.Errorf("can't declare queue: %w", err)
	}

	err = r.ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		},
	)

	if err != nil {
		return fmt.Errorf("can't publish message: %w", err)
	}
	return nil
}

func (r *RabbitMQ) Consume(ctx context.Context, f func(m []byte) error) error {
	q, err := r.ch.QueueDeclare(
		r.config.GetQueue(), // name of the queue
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)

	if err != nil {
		return fmt.Errorf("Failed to declare a RabbitMQ queue: %w", err)
	}

	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err != nil {
		return fmt.Errorf("Failed to register a RabbitMQ consumer: %w", err)
	}

	msgsCount := q.Messages
	if msgsCount == 0 {
		return nil
	}

	for msg := range msgs {
		err := f(msg.Body)
		if err != nil {
			fmt.Println("can't process message")
		}
		msg.Ack(true)
	}

	return nil
}
