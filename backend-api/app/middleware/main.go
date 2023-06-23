package middleware

import (
    "log"
    "time"
    "strconv"
    "context"
    "github.com/gofiber/fiber/v2"
    "github.com/anuragpal/city-falcon-test/api/app/utils"
    "github.com/anuragpal/city-falcon-test/api/app/models"
)


type M struct {

}

func (m *M) Interceptor(c *fiber.Ctx) error {
    if c.Method() == "GET" {
        var ctx = context.Background()
        url := c.OriginalURL()
        h := utils.Md5Hash(url)
        val, err := models.RC.Get(ctx, h).Result()
        if err != nil {
            err := models.RC.Set(ctx, h, "1", 24 * 60 * 60 * time.Second).Err()
            if err != nil {
                log.Println(err)
            }
            c.Locals("key", h)
            c.Locals("count", "1")
            return c.Next()
        }

        v, _ := strconv.Atoi(val)
        err = models.RC.Set(ctx, h, (v + 1), 24 * 60 * 60 * time.Second).Err()
        if err != nil {
            log.Println(err)
        }
        c.Locals("key", h)
        c.Locals("count", strconv.Itoa(v + 1))
    }
    return c.Next()
}