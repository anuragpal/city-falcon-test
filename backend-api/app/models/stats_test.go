package models

import (
    "testing"
    "github.com/gofiber/fiber/v2"
)

func TestPSqlStats(t *testing.T) {
    ps := PgStats{}
    params := ListParams{
        RecordPerPage: 10,
        Page: 1,
    }

    result, status := ps.PSqlStats(params)
    if status != fiber.StatusOK {
        t.Errorf("Expected status code 200, got %d", status)
    }

    data := result["data"].([]Stats)
    if len(data) != 10 {
        t.Errorf("Expected 1 row in result, got %d", len(data))
    }
}