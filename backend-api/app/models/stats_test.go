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
        Query: "test-not-exists",
    }

    result, status := ps.PSqlStats(params)
    if status != fiber.StatusOK {
        t.Errorf("Expected status code 200, got %d", status)
    }

    data := result["data"].([]Stats)
    if len(data) != 0 {
        t.Errorf("Expected 0 row in result, got %d", len(data))
    }

    params = ListParams{
        RecordPerPage: 0,
        Page: 0,
        Query: "test-not-exists",
    }

    result, status = ps.PSqlStats(params)
    if status != fiber.StatusOK {
        t.Errorf("Expected status code 200, got %d", status)
    }

    data = result["data"].([]Stats)
    if len(data) != 0 {
        t.Errorf("Expected 0 row in result, got %d", len(data))
    }

    params = ListParams{
        RecordPerPage: 0,
        Page: 0,
        Query: "",
    }

    result, status = ps.PSqlStats(params)
    if status != fiber.StatusOK {
        t.Errorf("Expected status code 200, got %d", status)
    }

    data = result["data"].([]Stats)
    if len(data) == 0 {
        t.Errorf("")
    }
}