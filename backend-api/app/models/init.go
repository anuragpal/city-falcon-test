package models

import (
    "os"
    "log"
    "strings"
    "database/sql"
    _ "github.com/lib/pq"
    "github.com/redis/go-redis/v9"
)

var (
    DB        *sql.DB
    RC        *redis.Client
    RC1       *redis.Client
)

func init() {
    var err error
    conn_string := []string{"postgres", ":", "//", os.Getenv("POSTGRES_USER"), ":", os.Getenv("POSTGRES_PASSWORD"), "@", os.Getenv("POSTGRES_URI"), "/", os.Getenv("POSTGRES_DB"), "?sslmode=disable"}
	conn := strings.Join(conn_string, "")

    DB, err = sql.Open("postgres", conn)
    if err != nil {
        log.Println(err)
    }

    RC = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URI"),
        Password: "",
        DB:       0,
    })

    RC1 = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_URI"),
        Password: "",
        DB:       0,
    })

    err = DB.Ping()
    if err != nil {
        log.Println(err)
    }
}