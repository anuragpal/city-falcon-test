package utils

import (
    "testing"
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