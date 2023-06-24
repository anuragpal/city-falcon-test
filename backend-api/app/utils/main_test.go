package utils

import (
    "time"
    "context"
    "testing"
    "github.com/anuragpal/city-falcon-test/api/app/models"
)

func TestMd5Hash(t *testing.T) {
    cases := []struct {
        str string
        want string
    }{
        {
            str: "testmd5hash",
            want: "8f096643c6e9c652907d8e493290055a",
        },
        {
            str: "",
            want: "",
        },
    }

    for _, tc := range cases {
        got := Md5Hash(tc.str)
        if got != tc.want {
            t.Errorf("Md5Hash(%s) = %s, want %s", tc.str, got, tc.want)
        }
    }
}

func TestRemoveCache(t *testing.T) {
    var ctx = context.Background()
    if err := models.RC.Set(ctx, "tstuser1-1", "value1", 24 * 60 * 60 * time.Second).Err(); err != nil {
        t.Errorf("Failed to put key-value pair in redis: %v", err)
    }

    if err := models.RC.Set(ctx, "tstuser2-2", "value2", 24 * 60 * 60 * time.Second).Err(); err != nil {
        t.Errorf("Failed to put key-value pair in redis: %v", err)
    }

    if err := RemoveCache("tstuser1"); err != nil {
        t.Errorf("Failed to remove cache: %v", err)
    }

    if _, err := models.RC.Get(ctx, "tstuser1-1").Result(); err == nil {
        t.Errorf("Expected key 'tstuser1-1' to be removed from cache")
    }

    if err := RemoveCache("tstuser2*"); err != nil {
        t.Errorf("Failed to remove cache: %v", err)
    }

    if _, err := models.RC.Get(ctx, "tstuser2-2").Result(); err == nil {
        t.Errorf("Expected key 'tstuser2-2' to be removed from cache")
    }
}