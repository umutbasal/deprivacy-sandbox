package main

import (
	_ "embed"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Deprivacy Sandbox - Time-based")
	})

	app.Listen(":8080")

}
