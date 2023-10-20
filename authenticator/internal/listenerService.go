package internal

import (
	"Projeect/utils/broker"
	"Projeect/utils/datasource"
	"github.com/sirupsen/logrus"
)

const (
	QUEUE = "reqs"
)

var (
	rabbitMQ *broker.RabbitMQ
	psql     *datasource.PSQL
)

func Run() {
	InitTools()

	go rabbitMQ.ListenForMessages()
	listener()
}

func InitTools() {
	var err error
	psql, err = datasource.InitializePSQL()
	if err != nil {
		logrus.Fatalf("Error initializing database: %v", err)
	}

	rabbitMQ, err = broker.InitRabbitMQ(QUEUE)
	if err != nil {
		logrus.Fatalf("Error initializing rabbitMQ: %v", err)
	}
	//defer func(mq *broker.RabbitMQ) {
	//	err := mq.Close()
	//	if err != nil {
	//		logrus.Errorf("Error closing rabbitMQ: %v", err)
	//	}
	//}(rabbitMQ)
}

func listener() {
	for msg := range rabbitMQ.CodeReqs {
		go processReq(msg)
	}
}
