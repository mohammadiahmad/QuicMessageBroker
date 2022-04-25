package http

import (
	"fmt"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Run(Host string, Port string) {
	app := fiber.New()

	app.Use(recover.New())

	prometheus := fiberprometheus.NewWith("http", "dispatching", "gabriel")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	address := fmt.Sprintf("%s:%s", Host, Port)
	app.Listen(address)
	fmt.Println("Http Service Started Successfully.")
}
