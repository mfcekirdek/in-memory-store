package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	model2 "gitlab.com/mfcekirdek/in-memory-store/pkg/model"
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

func ParseBody(body []byte) *model2.BaseResponse {
	var actualResponseBody *model2.BaseResponse
	err := json.Unmarshal(body, &actualResponseBody)
	if err != nil {
		fmt.Println(err)
	}
	return actualResponseBody
}
