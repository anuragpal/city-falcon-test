package middleware

import (
    "github.com/gofiber/fiber/v2"
)


type M struct {

}

func (m *M) Interceptor(c *fiber.Ctx) error {
    return c.Next()
}