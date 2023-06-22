package stats

import (
    "strconv"
    "github.com/gofiber/fiber/v2"
    "github.com/anuragpal/city-falcon-test/api/app/models"
)

func PsqlSlowQuery(c *fiber.Ctx) error {
    var params      models.ListParams
    ps := models.PgStats{}
    params.Query = c.Query("q")
    params.SortBy = c.Query("sort_by")
    params.OrderBy = c.Query("order_by")
    params.Page, _ = strconv.ParseInt(c.Query("page"), 10, 64)
    params.RecordPerPage, _ = strconv.ParseInt(c.Query("rpp"), 10, 64)

    result, status := ps.PSqlStats(params)
    return c.Status(status).JSON(result)
}