package models

import (
    "os"
    "fmt"
    "strings"
    "database/sql"
    _ "github.com/lib/pq"
)

var (
    DB        *sql.DB
)

func init() {
    var err error
    conn_string := []string{"postgres", ":", "//", os.Getenv("POSTGRES_USER"), ":", os.Getenv("POSTGRES_PASSWORD"), "@", os.Getenv("POSTGRES_URI"), "/", os.Getenv("POSTGRES_DB"), "?sslmode=disable"}
	conn := strings.Join(conn_string, "")

    DB, err = sql.Open("postgres", conn)
    if err != nil {
        fmt.Println("DB Error", err)
    } else {
        fmt.Println("DB Connected anurag")
    }
}