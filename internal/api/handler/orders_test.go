package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_orders_getOrders(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	o := &order{l: logger.Sugar()}

	type args struct {
		w           *httptest.ResponseRecorder
		method      string
		contentType string
		url         string
		body        io.Reader
	}
	type want struct {
		code        int
		contentType string
		response    string
		err         bool
	}
	tests := []struct {
		name string
		o    *order
		args args
		want want
	}{
		{
			name: "Test case #1",
			o:    o,
			args: args{
				w:           httptest.NewRecorder(),
				method:      http.MethodGet,
				url:         "/api/user/orders",
				contentType: "application/json",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			r.Header.Set("Content-Type", tt.args.contentType)

			tt.o.getOrders(tt.args.w, r)
			assert.Equal(t, tt.want.code, tt.args.w.Code)
			assert.Equal(t, tt.want.contentType, tt.args.w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, tt.args.w.Body.String())
			}
		})
	}
}

func Test_orders_addBonus(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	o := &order{l: logger.Sugar()}

	type args struct {
		w           *httptest.ResponseRecorder
		method      string
		contentType string
		url         string
		body        io.Reader
	}
	type want struct {
		code        int
		contentType string
		response    string
		err         bool
	}
	tests := []struct {
		name string
		o    *order
		args args
		want want
	}{
		{
			name: "Test case #1",
			o:    o,
			args: args{
				w:           httptest.NewRecorder(),
				method:      http.MethodPost,
				url:         "/api/user/orders",
				contentType: "application/json",
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			r.Header.Set("Content-Type", tt.args.contentType)

			tt.o.addBonus(tt.args.w, r)
			assert.Equal(t, tt.want.code, tt.args.w.Code)
			assert.Equal(t, tt.want.contentType, tt.args.w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, tt.args.w.Body.String())
			}
		})
	}
}
