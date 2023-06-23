package app

import (
    "testing"
    "net/http"
    "io/ioutil"
    "github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
    i := App{}
    app := i.Initialize()

    req, _ := http.NewRequest(
        "GET",
        "/",
        nil,
    )

    res, err := app.Test(req, -1)
    assert.Equalf(t, false, err != nil, "index route")
    assert.Equalf(t, 200, res.StatusCode, "index route")
    body, err := ioutil.ReadAll(res.Body)
    assert.Nilf(t, err, "index route")
    assert.Equalf(t, "OK", string(body), "index route")
}

func TestStats(t *testing.T) {
    i := App{}
    app := i.Initialize()

    req, _ := http.NewRequest(
        "GET",
        "/v1/psq",
        nil,
    )

    res, err := app.Test(req, -1)
    assert.Equalf(t, false, err != nil, "index route")
    assert.Equalf(t, 200, res.StatusCode, "index route")
}