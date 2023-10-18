package handlers

import (
	"Projeect/utils/broker"
	"Projeect/utils/datasource"
	"github.com/sirupsen/logrus"
)

var (
	rabbitMQ *broker.RabbitMQ
	psql     *datasource.PSQL
)

func InitTools() {
	var err error
	psql, err = datasource.InitializePSQL()
	if err != nil {
		logrus.Fatalf("Error initializing database: %v", err)
	}

	rabbitMQ, err = broker.InitRabbitMQ()
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
