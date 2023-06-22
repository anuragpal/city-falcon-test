package users

import (
    "fmt"
    "strconv"
    // "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/anuragpal/city-falcon-test/api/app/models"
)

func Add(c *fiber.Ctx) error {
    u := models.User{}
    if err := c.BodyParser(&u); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
    fmt.Println(u.LastName)
    result, status := u.Add()
    return c.Status(status).JSON(result)
}

func Update(c *fiber.Ctx) error {
    u := models.User{}
    if err := c.BodyParser(&u); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
    
    u.Id = c.Params("id")
    result, status := u.Update()
    return c.Status(status).JSON(result)
}

func List(c *fiber.Ctx) error {
    var params      models.ListParams
    u := models.User{}
    params.Query = c.Query("q")
    params.SortBy = c.Query("sort_by")
    params.OrderBy = c.Query("order_by")
    params.Page, _ = strconv.ParseInt(c.Query("page"), 10, 64)
    params.RecordPerPage, _ = strconv.ParseInt(c.Query("rpp"), 10, 64)

    result, status := u.List(params)
    return c.Status(status).JSON(result)
}

func Delete(c *fiber.Ctx) error {
    u := models.User{}
    u.Id = c.Params("id")
    result, status := u.Remove()
    return c.Status(status).JSON(result)
}

func Details(c *fiber.Ctx) error {
    u := models.User{}
    u.Id = c.Params("id")
    result, status := u.Details()
    return c.Status(status).JSON(result)
}