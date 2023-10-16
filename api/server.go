package api

import (
	"Projeect/utils/datasource"
	"github.com/sirupsen/logrus"
)

func run() error {
	psql, err := datasource.InitializePSQL()
	if err != nil {
		logrus.Fatalf("Error initializing database: %v", err)
	}
}
