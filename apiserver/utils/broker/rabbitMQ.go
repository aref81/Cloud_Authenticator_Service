package broker

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

var (
	rabbitURL = os.Getenv("RB_URL")
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func InitRabbitMQ() (mq *RabbitMQ, err error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}

func (mq *RabbitMQ) Close() error {
	err := mq.channel.Close()
	if err != nil {
		return err
	}

	err = mq.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (mq *RabbitMQ) Publish(ctx context.Context, queue, body string) error {
	publishing := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}

	err := mq.channel.PublishWithContext(ctx, "", queue, false, false, publishing)
	if err != nil {
		return err
	}

	return nil
}
