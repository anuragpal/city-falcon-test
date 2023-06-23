package utils

import (
    "context"
    "crypto/md5"
    "encoding/hex"
    "github.com/anuragpal/city-falcon-test/api/app/models"
)

func Md5Hash(str string) string {
    if str == "" {
        return ""
    }
    hash := md5.Sum([]byte(str))
    return hex.EncodeToString(hash[:])
}

func RemoveCache(prefix string) error {
    var ctx = context.Background()
    iter := models.RC.Scan(ctx, 0, "users*", 0).Iterator()
    for iter.Next(ctx) {
        key := iter.Val()
        if err := models.RC.Del(ctx, key).Err(); err != nil {
            return err
        }
    }
    if err := iter.Err(); err != nil {
        return err
    }
    return nil
}