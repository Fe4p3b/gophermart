package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_auth_register(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	a := &auth{l: logger.Sugar()}

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
		a    *auth
		args args
		want want
	}{
		{
			name: "Test case #1",
			a:    a,
			args: args{
				w:           httptest.NewRecorder(),
				method:      http.MethodPost,
				url:         "/api/user/register",
				contentType: "application/json",
				body: strings.NewReader(`{
					"login": "user",
					"password": "password"
				} `),
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

			tt.a.register(tt.args.w, r)
			assert.Equal(t, tt.want.code, tt.args.w.Code)
			assert.Equal(t, tt.want.contentType, tt.args.w.Header().Get("Content-Type"))
			if tt.want.response != "" {
				assert.Equal(t, tt.want.response, tt.args.w.Body.String())
			}
		})
	}
}

func Test_auth_login(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	a := &auth{l: logger.Sugar()}

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
		a    *auth
		args args
		want want
	}{
		{
			name: "Test case #1",
			a:    a,
			args: args{
				w:           httptest.NewRecorder(),
				method:      http.MethodPost,
				url:         "/api/user/login",
				contentType: "application/json",
				body: strings.NewReader(`{
					"login": "user",
					"password": "password"
				} `),
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
			},
		},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(tt.args.method, tt.args.url, tt.args.body)
		r.Header.Set("Content-Type", tt.args.contentType)

		tt.a.login(tt.args.w, r)
		assert.Equal(t, tt.want.code, tt.args.w.Code)
		assert.Equal(t, tt.want.contentType, tt.args.w.Header().Get("Content-Type"))
		if tt.want.response != "" {
			assert.Equal(t, tt.want.response, tt.args.w.Body.String())
		}
	}
}
