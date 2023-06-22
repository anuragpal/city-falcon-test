package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/anuragpal/city-falcon-test/api/app"
)

func main() {
    app := Setup()
    app.Listen(":3000")
}

func Setup() *fiber.App {
    instance := app.App{}
    return instance.Initialize()
}