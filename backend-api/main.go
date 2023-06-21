package main

import (
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := Setup()
    app.Listen(":3000")
}

func Setup() *fiber.App {
    app := fiber.New()
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("OK")
    })
    return app
}