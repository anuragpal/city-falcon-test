package app

import (
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/anuragpal/city-falcon-test/api/app/users"
    "github.com/anuragpal/city-falcon-test/api/app/stats"
    "github.com/anuragpal/city-falcon-test/api/app/middleware"
)

type App struct {

}

func (a *App) Initialize() *fiber.App {
    app := fiber.New()
    app.Use(logger.New(logger.Config{
        Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
    }))
    
    app.Get("/", Default)

    m := middleware.M{}

    // managed version based API's
    v1 := app.Group("/v1", m.Interceptor)

    // User's crud operation
    v1.Get("/users", users.List)
    v1.Get("/user/:id", users.Details)
    v1.Post("/user", users.Add)
    v1.Put("/user/:id", users.Update)
    v1.Delete("/user/:id", users.Delete)

    // Psql Slow query data retrival
    v1.Get("/psq", stats.PsqlSlowQuery)

    return app
}