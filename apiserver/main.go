package main

import (
	"Projeect/api"
	"github.com/sirupsen/logrus"
)

func main() {
	err := api.Run()
	if err != nil {
		logrus.Fatalf(err.Error())
		return
	}
}
