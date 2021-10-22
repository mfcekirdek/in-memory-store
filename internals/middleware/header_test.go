//+build unit

package middleware

import (
	"github.com/golang/mock/gomock"
	"gitlab.com/mfcekirdek/in-memory-store/mocks"
	"gitlab.com/mfcekirdek/in-memory-store/test"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHeaderMiddleware_ServeHTTP(t *testing.T) {
	mockController := gomock.NewController(t)
	mockHandler := mocks.NewMockStoreHandler(mockController)
	mockHandler.
		EXPECT().
		ServeHTTP(gomock.Any(), gomock.Any()).
		Return().
		MaxTimes(1)

	w, r := test.CreateHTTPReq(http.MethodGet, "/api/v1/store", nil)
	wOptions, rOptions := test.CreateHTTPReq(http.MethodOptions, "/api/v1/store", nil)

	var wantHeaders http.Header = map[string][]string{
		"Content-Type":                 {"application/json"},
		"Access-Control-Allow-Origin":  {"*"},
		"Access-Control-Allow-Methods": {"GET, OPTIONS, PUT, DELETE"},
		"Access-Control-Allow-Headers": {"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
	}
	type fields struct {
		handler http.Handler
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   http.Header
	}{
		{"sets headers", fields{handler: mockHandler}, args{
			w: w,
			r: r,
		}, wantHeaders},
		{"options request", fields{handler: mockHandler}, args{
			w: wOptions,
			r: rOptions,
		}, wantHeaders},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &HeaderMiddleware{
				handler: tt.fields.handler,
			}

			l.ServeHTTP(tt.args.w, tt.args.r)
			if got := tt.args.w.Header(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHeaderMiddleware(t *testing.T) {
	mockController := gomock.NewController(t)
	mockHandler := mocks.NewMockStoreHandler(mockController)

	type args struct {
		handlerToWrap http.Handler
	}
	tests := []struct {
		name string
		args args
		want *HeaderMiddleware
	}{
		{"new instance", args{handlerToWrap: mockHandler}, &HeaderMiddleware{handler: mockHandler}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHeaderMiddleware(tt.args.handlerToWrap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHeaderMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
