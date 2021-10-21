//+build unit

package handler

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/mfcekirdek/in-memory-store/internals/service"
	"gitlab.com/mfcekirdek/in-memory-store/mocks"
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

//func Test_storeHandler_Flush(t *testing.T) {
//	mockController := gomock.NewController(t)
//	mockService := mocks.NewMockStoreService(mockController)
//	mockService.
//		EXPECT().
//		Flush().
//		Return(map[string]string{}).
//		Times(1)
//
//	w, r := test.CreateHTTPReq(http.MethodDelete, "/api/v1/store", nil)
//
//	type fields struct {
//		service service.StoreService
//	}
//	type args struct {
//		w *httptest.ResponseRecorder
//		r *http.Request
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   *model.BaseResponse
//	}{
//		{"[200] - flush all store data", fields{service: mockService}, args{
//			w: w,
//			r: r,
//		}, &model.BaseResponse{
//			Data:        map[string]interface{}{},
//			Description: "all items deleted",
//		}},
//		{"[405] - StatusMethodNotAllowed", fields{service: mockService}, args{
//			w: w,
//			r: r,
//		}, &model.BaseResponse{
//			Data:        map[string]interface{}{},
//			Description: "method not allowed",
//		}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &storeHandler{
//				service: tt.fields.service,
//			}
//
//			s.Flush(tt.args.w, tt.args.r)
//			if got := test.ParseBody(w.Body.Bytes()); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewStoreService() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

//func Test_storeHandler_ServeHTTP(t *testing.T) {
//	type fields struct {
//		service service.StoreService
//	}
//	type args struct {
//		w http.ResponseWriter
//		r *http.Request
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := &storeHandler{
//				service: tt.fields.service,
//			}
//		})
//	}
//}
