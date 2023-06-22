package models

import (
    "strings"
    "strconv"
    "github.com/gofiber/fiber/v2"
)

type PgStats struct {
    
}

type Stats struct {
    Query           string             `json:"query,omitempty"`
    ExecutedTime    float64            `json:"executed_time,omitempty"`
    Calls           int                `json:"calls,omitempty"`
    Mean            float64            `json:"mean,omitempty"`
    CPU             float64            `json:"cpu,omitempty"`
}

func (ps PgStats) PSqlStats(l ListParams) (fiber.Map, int) {
    var (
        condition_count             int
        params                      []interface{}
        conditions                  []string
        total                        int
    )

    conditions = []string{}
    condition_count++
    params = append(params, 0)
    conditions = append(conditions, " p.total_time > $"+strconv.Itoa(condition_count))

    if l.Query != "" {
        condition_count = condition_count + 1
        params = append(params, "%"+l.Query+"%")
        conditions = append(conditions, " (p.query LIKE $1)")
    }

    stats := []Stats{}
    if l.RecordPerPage <= 0 {
        l.RecordPerPage = 10
    }

    if l.Page <= 0  {
        l.Page = 1
    }

    offset := (l.Page - 1) * l.RecordPerPage

    stmt, err := DB.Prepare(`SELECT count(p.queryid) as cnt FROM pg_stat_statements as p WHERE ` + strings.Join(conditions, " AND "))
    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    err = stmt.QueryRow(params...).Scan(&total)
    if err != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    limit_param_count := condition_count + 1
    params = append(params, l.RecordPerPage)

    offset_param_count := limit_param_count + 1
    params = append(params, offset)

    sql := `
    SELECT
        substring(query, 1, 50) AS short_query,
        round(total_time::numeric, 2) AS total_exec_time,
        calls,
        round(mean_time::numeric, 2) AS mean,
        round((100 * total_time /
        sum(total_time::numeric) OVER ())::numeric, 2) AS percentage_cpu
    FROM pg_stat_statements as p
    WHERE ` + strings.Join(conditions, " AND ") + " LIMIT $"+strconv.Itoa(limit_param_count) + " OFFSET $" + strconv.Itoa(offset_param_count)
    rows, errRows := DB.Query(sql, params...)

    if errRows != nil {
        return fiber.Map{}, fiber.StatusInternalServerError
    }

    for rows.Next() {
        stat := Stats{}
        errScan := rows.Scan(&stat.Query, &stat.ExecutedTime, &stat.Calls, &stat.Mean, &stat.CPU)
        if errScan != nil {
            return fiber.Map{}, fiber.StatusInternalServerError
        }
        stats = append(stats, stat)
    }
    
    return fiber.Map{
        "data": stats,
        "rpp": l.RecordPerPage,
        "sort_by": l.SortBy,
        "order_by": l.OrderBy,
        "page": l.Page,
        "total": total,
    }, fiber.StatusOK
}