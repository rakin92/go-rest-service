package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakin92/go-rest-service/pkg/logger"
)

func Test_jwtFromHeader(t *testing.T) {
	hed := http.Header{}
	req := http.Request{Header: hed}
	ctx := gin.Context{Request: &req}

	const authHeader = "Authorization"
	type args struct {
		c   *gin.Context
		key string
		tok string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "jwtFromHeader missing key",
			args: args{
				c:   &ctx,
				key: "",
				tok: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "jwtFromHeader bad token",
			args: args{
				c:   &ctx,
				key: authHeader,
				tok: "BadToken",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "jwtFromHeader good token",
			args: args{
				c:   &ctx,
				key: authHeader,
				tok: "Bearer Good",
			},
			want:    "Good",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req.Header.Set(authHeader, tt.args.tok)

			got, err := jwtFromHeader(tt.args.c, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("jwtFromHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("jwtFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_apiKeyFromHeader(t *testing.T) {
	hed := http.Header{}
	req := http.Request{Header: hed}
	ctx := gin.Context{Request: &req}

	type args struct {
		c   *gin.Context
		key string
		tok string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid token",
			args: args{
				c:   &ctx,
				key: APIKeyHeader,
				tok: "a_valid_api_kwy",
			},
			want:    "a_valid_api_kwy",
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				c:   &ctx,
				key: "invalid_key",
				tok: "",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "invalid token",
			args: args{
				c:   &ctx,
				key: APIKeyHeader,
				tok: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req.Header.Set(APIKeyHeader, tt.args.tok)
			got, err := apiKeyFromHeader(tt.args.c, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("apiKeyFromHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("apiKeyFromHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tokenFromQuery(t *testing.T) {
	hed := http.Header{}
	url := url.URL{}
	req := http.Request{Header: hed, URL: &url}
	ctx := gin.Context{Request: &req}

	type args struct {
		c   *gin.Context
		key string
		tok string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid token",
			args: args{
				c:   &ctx,
				key: "token",
				tok: "good_token",
			},
			want:    "good_token",
			wantErr: false,
		},
		{
			name: "invalid api key",
			args: args{
				c:   &ctx,
				key: "bad_token",
				tok: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := fmt.Sprintf("%v=%v", tt.args.key, tt.args.tok)
			req.URL.RawQuery = q
			got, err := tokenFromQuery(tt.args.c, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenFromQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tokenFromQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tokenFromCookie(t *testing.T) {
	hed := http.Header{}
	req := http.Request{Header: hed}
	ctx := gin.Context{Request: &req}

	type args struct {
		c   *gin.Context
		key string
		tok string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid token",
			args: args{
				c:   &ctx,
				key: "jwt",
				tok: "good_token",
			},
			want:    "good_token",
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				c:   &ctx,
				key: "bad_jwt",
				tok: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cookie := &http.Cookie{
				Name:   tt.args.key,
				Value:  tt.args.tok,
				MaxAge: 300,
			}
			req.AddCookie(cookie)
			logger.Info("test")
			got, err := tokenFromCookie(tt.args.c, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenFromCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tokenFromCookie() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tokenFromParam(t *testing.T) {
	hed := http.Header{}
	req := http.Request{Header: hed}
	ctx := gin.Context{Request: &req}

	type args struct {
		c   *gin.Context
		key string
		tok string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "valid token",
			args: args{
				c:   &ctx,
				key: "token",
				tok: "good_token",
			},
			want:    "good_token",
			wantErr: false,
		},
		{
			name: "invalid key",
			args: args{
				c:   &ctx,
				key: "bad_token",
				tok: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.c.Params = []gin.Param{
				{
					Key:   tt.args.key,
					Value: tt.args.tok,
				},
			}
			got, err := tokenFromParam(tt.args.c, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokenFromParam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tokenFromParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
