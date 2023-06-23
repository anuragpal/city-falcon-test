package users

import (
    "os"
    "log"
    "time"
    "context"
    "strings"
    "strconv"
    "encoding/json"
    "github.com/gofiber/fiber/v2"
    "github.com/anuragpal/city-falcon-test/api/app/utils"
    "github.com/anuragpal/city-falcon-test/api/app/models"
)

type CachedData struct {
    Result              fiber.Map               `json:"result"`
    Status              int                     `json:"status"`
}

func Add(c *fiber.Ctx) error {
    u := models.User{}
    if err := c.BodyParser(&u); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
    
    err := utils.RemoveCache("users_")
    if err != nil {
        log.Println("Failed to remove cache.")
    }
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
    err := utils.RemoveCache("users_")
    if err != nil {
        log.Println("Failed to remove cache.")
    }
    result, status := u.Update()
    return c.Status(status).JSON(result)
}

func List(c *fiber.Ctx) error {
    h := c.Locals("key").(string)
    cnt, _ := strconv.Atoi(c.Locals("count").(string))
    CACHE_MIN_HITS, _ := strconv.Atoi(os.Getenv("CACHE_MIN_HITS"))

    var (
        params      models.ListParams
        result      fiber.Map
        status      int
        data        CachedData
    )
    var ctx = context.Background()

    u := models.User{}
    params.Query = c.Query("q")
    params.SortBy = c.Query("sort_by")
    params.OrderBy = c.Query("order_by")
    params.Page, _ = strconv.ParseInt(c.Query("page"), 10, 64)
    params.RecordPerPage, _ = strconv.ParseInt(c.Query("rpp"), 10, 64)

    if cnt >= CACHE_MIN_HITS {
        cmb := []string{"users", h}
        hKey := strings.Join(cmb, "_")
        val, err := models.RC.Get(ctx, hKey).Result()
        if err != nil {
            result, status = u.List(params)
            data.Result = result
            data.Status = status
            cData, errData := json.Marshal(data)
            if errData != nil {
                log.Println(errData)
            }

            err = models.RC.Set(ctx, hKey, cData, 24 * 60 * 60 * time.Second).Err()
            if err != nil {
                log.Println(err)
            }
            return c.Status(status).JSON(result)
        }
        
        err = json.Unmarshal([]byte(val), &data)
        if err != nil {
            log.Println(err)
        }
        return c.Status(data.Status).JSON(data.Result)
    } else {
        result, status = u.List(params)
    }
    
    return c.Status(status).JSON(result)
}

func Delete(c *fiber.Ctx) error {
    u := models.User{}
    u.Id = c.Params("id")
    err := utils.RemoveCache("users_")
    if err != nil {
        log.Println("Failed to remove cache.")
    }
    result, status := u.Remove()
    return c.Status(status).JSON(result)
}

func Details(c *fiber.Ctx) error {
    u := models.User{}
    u.Id = c.Params("id")
    result, status := u.Details()
    return c.Status(status).JSON(result)
}