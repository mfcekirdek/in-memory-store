package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gitlab.com/mfcekirdek/in-memory-store/internals/model"
	"io"
	"net/http"
	"net/http/httptest"
)

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

func ParseBody(body []byte) *model.BaseResponse {
	var actualResponseBody *model.BaseResponse
	err := json.Unmarshal(body, &actualResponseBody)
	if err != nil {
		fmt.Println(err)
	}
	return actualResponseBody
}
