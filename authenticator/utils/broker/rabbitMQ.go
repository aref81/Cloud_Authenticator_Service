package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	rabbitURL = os.Getenv("RB_URL")
)

type RabbitMQ struct {
	CodeReqs chan string
	conn     *amqp.Connection
	channel  *amqp.Channel
	msgs     <-chan amqp.Delivery
}

func InitRabbitMQ(queue string) (mq *RabbitMQ, err error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		queue,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:     conn,
		channel:  channel,
		msgs:     msgs,
		CodeReqs: make(chan string, 10),
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

func (mq *RabbitMQ) ListenForMessages() {
	logrus.Infof("Listening for messages from rabbitMQ")
	for msg := range mq.msgs {
		logrus.Infof("Received a message: %s", msg.MessageId)
		mq.CodeReqs <- string(msg.Body)
	}
}
