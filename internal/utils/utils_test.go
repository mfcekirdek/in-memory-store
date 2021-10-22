//+build unit

package utils

import (
	"gitlab.com/mfcekirdek/in-memory-store/internal/model"
	"gitlab.com/mfcekirdek/in-memory-store/test"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGenerateResponse(t *testing.T) {
	data := map[string]string{"foo": "bar"}
	description := "item created"
	type args struct {
		data        interface{}
		description string
	}
	tests := []struct {
		name string
		args args
		want *model.BaseResponse
	}{
		{"Generate Base Response", args{
			data:        data,
			description: description,
		}, &model.BaseResponse{
			Data:        data,
			Description: description,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateResponse(tt.args.data, tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandleError(t *testing.T) {
	w404, r404 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/foo", nil)
	w400, r400 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/foo", nil)
	w405, r405 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/foo", nil)

	type args struct {
		w      *httptest.ResponseRecorder
		r      *http.Request
		status int
	}
	tests := []struct {
		name           string
		args           args
		want           *model.BaseResponse
		wantStatusCode int
	}{
		{"not found", args{
			w:      w404,
			r:      r404,
			status: http.StatusNotFound,
		}, &model.BaseResponse{
			Data:        nil,
			Description: "not found",
		}, http.StatusNotFound},
		{"bad input parameter/body", args{
			w:      w400,
			r:      r400,
			status: http.StatusBadRequest,
		}, &model.BaseResponse{
			Data:        nil,
			Description: "bad input parameter/body",
		}, http.StatusBadRequest},
		{"method not allowed", args{
			w:      w405,
			r:      r405,
			status: http.StatusMethodNotAllowed,
		}, &model.BaseResponse{
			Data:        nil,
			Description: "method not allowed",
		}, http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleError(tt.args.w, tt.args.r, tt.args.status)
			if got := test.ParseBody(tt.args.w.Body.Bytes()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HandleError() = %v, want %v", got, tt.want)
			}

			if statusCode := tt.args.w.Code; !reflect.DeepEqual(statusCode, tt.wantStatusCode) {
				t.Errorf("StatusCode = %v, want %v", statusCode, tt.wantStatusCode)
			}
		})
	}
}

func TestReturnJSONResponse(t *testing.T) {
	w500, r500 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/foo", nil)
	w200, r200 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/foo", nil)

	type args struct {
		w           *httptest.ResponseRecorder
		r           *http.Request
		result      interface{}
		description string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
	}{
		{"[200] - success", args{
			w:           w200,
			r:           r200,
			result:      map[string]string{"foo": "bar"},
			description: "result",
		}, http.StatusOK},
		{"[500] - marshal error", args{
			w:           w500,
			r:           r500,
			result:      math.Inf(1),
			description: "test",
		}, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReturnJSONResponse(tt.args.w, tt.args.r, tt.args.result, tt.args.description)
			if statusCode := tt.args.w.Code; !reflect.DeepEqual(statusCode, tt.wantStatusCode) {
				t.Errorf("StatusCode = %v, want %v", statusCode, tt.wantStatusCode)
			}
		})
	}
}
