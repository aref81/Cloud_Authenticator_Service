package api

import (
	"Projeect/api/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

const (
	HttpServerPort = ":8000"
)

func Run() error {
	handlers.InitTools()

	e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/register", handlers.RegisterHandler)
	e.GET("/status", handlers.StatusHandler)

	err := e.Start(HttpServerPort)
	if err != nil {
		logrus.Errorf("Error starting http server: %v", err)
	}
	logrus.Infof("Server listening on port %s", HttpServerPort)

	return nil
}
