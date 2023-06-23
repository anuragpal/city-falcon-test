package utils

import (
    "crypto/md5"
    "encoding/hex"
)

func Md5Hash(str string) string {
    if str == "" {
        return ""
    }
    hash := md5.Sum([]byte(str))
    return hex.EncodeToString(hash[:])
}