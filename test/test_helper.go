// Package test provides helper functions to the test suites.
package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"gitlab.com/mfcekirdek/in-memory-store/pkg/model"
)

// Takes HTTP method type, endpoint and body paremeters.
// Creates a new HTTP Request and returns <httptest.ResponseRecorder> and <http.Request>.
func CreateHTTPReq(method, endpoint string, body io.Reader) (*httptest.ResponseRecorder, *http.Request) {
	if body == nil {
		const size = 512
		body = bytes.NewBuffer(make([]byte, size))
	}
	req := httptest.NewRequest(method, endpoint, body)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	rec := httptest.NewRecorder()
	return rec, req
}

// Parses []byte HTTP response body and returs as BaseResponse model.
func ParseBody(body []byte) *model.BaseResponse {
	var actualResponseBody *model.BaseResponse
	err := json.Unmarshal(body, &actualResponseBody)
	if err != nil {
		fmt.Println(err)
	}
	return actualResponseBody
}
