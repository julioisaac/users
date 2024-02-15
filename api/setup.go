package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/julioisaac/users/config"
	"github.com/julioisaac/users/logger"
	"os"
	"os/signal"
)

func SetupServer(port string, srv *fiber.App) error {
	logger.Logger.Infof("[INFO] App will initialize with %s env", config.GetString("APP_ENV"))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c //nolint:gosimple
		fmt.Println("Gracefully shutting down...")
		_ = srv.Shutdown()
	}()

	certFile := "certs/localhost.crt"
	keyFile := "certs/localhost.key"

	if err := srv.ListenTLS(":443", certFile, keyFile); err != nil {
		logger.Logger.Panic(err)
	}

	return nil
}

func CreateApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  config.GetDuration("HTTP_SERVER_READ_TIMEOUT_SECONDS"),
		WriteTimeout: config.GetDuration("HTTP_SERVER_WRITE_TIMEOUT_SECONDS"),
	})

	app.Use(recover.New())
	app.Use(compress.New())

	return app
}
