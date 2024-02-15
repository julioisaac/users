package api

import (
	"github.com/erkanzileli/nrfiber"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/julioisaac/users/api/health-api"
	"github.com/julioisaac/users/api/user"
	"github.com/julioisaac/users/config"
	"github.com/julioisaac/users/logger"
	"github.com/newrelic/go-agent/v3/newrelic"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title Users API
// @version 1.0
// @description Users Service API

// StartAPI server
func StartAPI() error {
	port := config.GetString("HTTP_PORT")

	app := CreateApp()

	nr, err := newrelic.NewApplication(
		newrelic.ConfigAppName("users"),
		newrelic.ConfigLicense("eu01xx4a8fa9b155456cedc7aa187016FFFFNRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		logger.Logger.Panic(err)
	}

	// Add the nrfiber middleware before other middlewares or routes
	app.Use(nrfiber.Middleware(nr))

	router := app.Group("/api")

	healthRouter := router.Group(healthAPI.HandlerPath)
	usersRouter := router.Group(usersAPI.HandlerPath)

	healthAPI.RegisterRoutes(healthRouter)
	usersAPI.RegisterRoutes(usersRouter)

	if config.GetBool("ENABLE_SWAGGER") {
		app.Get("/", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/docs/swagger/index.html", http.StatusMovedPermanently)
		}))
		app.Get("/docs/swagger/swagger.json", adaptor.HTTPHandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "./docs/swagger/swagger.json")
		}))
		app.Get("/docs/swagger/*", adaptor.HTTPHandlerFunc(httpSwagger.Handler(
			httpSwagger.URL("/docs/swagger/swagger.json"),
		)))
	}

	return SetupServer(port, app)
}
