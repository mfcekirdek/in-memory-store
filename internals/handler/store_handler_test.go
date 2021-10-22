//+build unit

package handler

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"gitlab.com/mfcekirdek/in-memory-store/internals/model"
	"gitlab.com/mfcekirdek/in-memory-store/internals/service"
	"gitlab.com/mfcekirdek/in-memory-store/internals/utils"
	"gitlab.com/mfcekirdek/in-memory-store/mocks"
	"gitlab.com/mfcekirdek/in-memory-store/test"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestNewStoreHandler(t *testing.T) {
	mockController := gomock.NewController(t)
	mockService := mocks.NewMockStoreService(mockController)

	type args struct {
		svc service.StoreService
	}
	tests := []struct {
		name string
		args args
		want StoreHandler
	}{
		{"New storeHandler instance", args{svc: mockService}, &storeHandler{service: mockService}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStoreHandler(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storeHandler_Flush(t *testing.T) {
	mockController := gomock.NewController(t)
	mockService := mocks.NewMockStoreService(mockController)
	mockService.
		EXPECT().
		Flush().
		Return(map[string]string{}).
		MaxTimes(1)

	wDelete200, rDelete200 := test.CreateHTTPReq(http.MethodDelete, "/api/v1/store", nil)
	wPost405, rPost405 := test.CreateHTTPReq(http.MethodPost, "/api/v1/store", nil)

	type fields struct {
		service service.StoreService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *model.BaseResponse
		wantStatusCode int
	}{
		{"[200] - flush all store data", fields{service: mockService}, args{
			w: wDelete200,
			r: rDelete200,
		}, &model.BaseResponse{
			Data:        map[string]interface{}{},
			Description: "all items deleted",
		}, http.StatusOK},
		{"[405] - StatusMethodNotAllowed", fields{}, args{
			w: wPost405,
			r: rPost405,
		}, &model.BaseResponse{
			Data:        nil,
			Description: "method not allowed",
		}, http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storeHandler{
				service: tt.fields.service,
			}

			s.Flush(tt.args.w, tt.args.r)
			if got := test.ParseBody(tt.args.w.Body.Bytes()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flush(() = %v, want %v", got, tt.want)
			}

			if statusCode := tt.args.w.Code; !reflect.DeepEqual(statusCode, tt.wantStatusCode) {
				t.Errorf("StatusCode = %v, want %v", statusCode, tt.wantStatusCode)
			}
		})
	}
}

func Test_storeHandler_ServeHTTP(t *testing.T) {
	mockController := gomock.NewController(t)
	mockService200 := mocks.NewMockStoreService(mockController)
	mockService200.
		EXPECT().
		Get("foo").
		Return(map[string]string{"foo": "bar"}, nil).
		MaxTimes(1)

	mockService200.
		EXPECT().
		Set("foo", "bar").
		Return(map[string]string{"foo": "bar"}, true).
		MaxTimes(1)

	mockService201 := mocks.NewMockStoreService(mockController)
	mockService201.
		EXPECT().
		Set("foo", "bar").
		Return(map[string]string{"foo": "bar"}, false).
		MaxTimes(1)

	mockService404 := mocks.NewMockStoreService(mockController)
	mockService404.
		EXPECT().
		Get("foo").
		Return(nil, utils.ErrNotFound).
		MaxTimes(1)

	wGet200, rGet200 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/foo", nil)
	wGet400, rGet400 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/", nil)
	wGet404, rGet404 := test.CreateHTTPReq(http.MethodGet, "/api/v1/store/foo", nil)
	wPost405, rPost405 := test.CreateHTTPReq(http.MethodPost, "/api/v1/store/foo", nil)

	requestBody, _ := json.Marshal(map[string]string{"value": "bar"})
	wPut200, rPut200 := test.CreateHTTPReq(http.MethodPut, "/api/v1/store/foo", bytes.NewBuffer(requestBody))
	wPut201, rPut201 := test.CreateHTTPReq(http.MethodPut, "/api/v1/store/foo", bytes.NewBuffer(requestBody))

	type fields struct {
		service service.StoreService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		want           *model.BaseResponse
		wantStatusCode int
	}{
		{"[200] GET - item fetched", fields{service: mockService200}, args{
			w: wGet200,
			r: rGet200,
		}, &model.BaseResponse{
			Data:        map[string]interface{}{"foo": "bar"},
			Description: "item fetched",
		}, http.StatusOK},
		{"[400] GET - bad input parameter/body", fields{}, args{
			w: wGet400,
			r: rGet400,
		}, &model.BaseResponse{
			Data:        nil,
			Description: "bad input parameter/body",
		}, http.StatusBadRequest},
		{"[404] GET - not found", fields{service: mockService404}, args{
			w: wGet404,
			r: rGet404,
		}, &model.BaseResponse{
			Data:        nil,
			Description: "not found",
		}, http.StatusNotFound},
		{"[405] - StatusMethodNotAllowed", fields{}, args{
			w: wPost405,
			r: rPost405,
		}, &model.BaseResponse{
			Data:        nil,
			Description: "method not allowed",
		}, http.StatusMethodNotAllowed},
		{"[200] PUT - item updated", fields{service: mockService200}, args{
			w: wPut200,
			r: rPut200,
		}, &model.BaseResponse{
			Data:        map[string]interface{}{"foo": "bar"},
			Description: "item updated",
		}, http.StatusOK},
		{"[201] PUT - item created", fields{service: mockService201}, args{
			w: wPut201,
			r: rPut201,
		}, &model.BaseResponse{
			Data:        map[string]interface{}{"foo": "bar"},
			Description: "item created",
		}, http.StatusCreated},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storeHandler{
				service: tt.fields.service,
			}

			s.ServeHTTP(tt.args.w, tt.args.r)

			if got := test.ParseBody(tt.args.w.Body.Bytes()); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServeHTTP() = %v, want %v", got, tt.want)
			}
			if statusCode := tt.args.w.Code; !reflect.DeepEqual(statusCode, tt.wantStatusCode) {
				t.Errorf("StatusCode = %v, want %v", statusCode, tt.wantStatusCode)
			}
		})
	}
}
