package main

import (
	"github.com/gofiber/compression"
	"github.com/gofiber/fiber"
)

func main() {

	app := fiber.New()

	ListenAddress := "0.0.0.0:2408"
	// Server Info
	app.Use(compression.New())
	app.Get("/", GwStat())
	app.All("/:apps/:key", GwWorker())

	app.Listen(ListenAddress)

}
