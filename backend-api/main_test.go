package main

import (
    "io/ioutil"
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestRoute(t *testing.T) {
    tests := []struct {
        description string

        // Test input
        route string

        // Expected output
        expectedError bool
        expectedCode  int
        expectedBody  string
    }{
        {
            description:   "index route",
            route:         "/",
            expectedError: false,
            expectedCode:  200,
            expectedBody:  "OK",
        },
        {
            description:   "non existing route",
            route:         "/i-dont-exist",
            expectedError: false,
            expectedCode:  404,
            expectedBody:  "Cannot GET /i-dont-exist",
        },
    }

    app := Setup()

    for _, test := range tests {
        req, _ := http.NewRequest(
            "GET",
            test.route,
            nil,
        )

        // Perform the request plain with the app.
        // The -1 disables request latency.
        res, err := app.Test(req, -1)

        // verify that no error occured, that is not expected
        assert.Equalf(t, test.expectedError, err != nil, test.description)

        // As expected errors lead to broken responses, the next
        // test case needs to be processed
        if test.expectedError {
            continue
        }

        // Verify if the status code is as expected
        assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

        // Read the response body
        body, err := ioutil.ReadAll(res.Body)

        // Reading the response body should work everytime, such that
        // the err variable should be nil
        assert.Nilf(t, err, test.description)

        // Verify, that the reponse body equals the expected body
        assert.Equalf(t, test.expectedBody, string(body), test.description)
    }
}