package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
)

func testRequest(h http.Handler, method string, address string, body string) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		address,
		bytes.NewBuffer([]byte(body)),
	)
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost || method == http.MethodPut {
		runes := []rune(body)
		if string(runes[0:1]) == "{" {
			req.Header.Add("Content-Type", "application/json")
		} else {
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	req.AddCookie(&http.Cookie{
		Name:  "hook_private",
		Value: "private",
	})

	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)

	return resp, nil
}
