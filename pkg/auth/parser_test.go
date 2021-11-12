package auth

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rakin92/go-rest-service/pkg/cfg"
	"github.com/rakin92/go-rest-service/pkg/consts"
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

func TestParseAPIKey(t *testing.T) {
	t.Run("api key from cookie", func(t *testing.T) {
		hed := http.Header{}
		url := &url.URL{}
		req := http.Request{Header: hed, URL: url}
		ctx := &gin.Context{Request: &req}
		svc := &cfg.Server{}
		cookie := &http.Cookie{
			Name:   "api_key",
			Value:  "token",
			MaxAge: 300,
		}
		req.AddCookie(cookie)
		gotApiKey, err := ParseAPIKey(ctx, svc)
		if (err != nil) != false {
			t.Errorf("ParseAPIKey() error = %v, wantErr %v", err, false)
			return
		}
		if gotApiKey != "token" {
			t.Errorf("ParseAPIKey() = %v, want %v", gotApiKey, "token")
		}
	})

	t.Run("api key from header", func(t *testing.T) {
		hed := http.Header{}
		url := &url.URL{}
		req := http.Request{Header: hed, URL: url}
		ctx := &gin.Context{Request: &req}
		svc := &cfg.Server{}
		req.Header.Set(APIKeyHeader, "token")
		gotApiKey, err := ParseAPIKey(ctx, svc)
		if (err != nil) != false {
			t.Errorf("ParseAPIKey() error = %v, wantErr %v", err, false)
			return
		}
		if gotApiKey != "token" {
			t.Errorf("ParseAPIKey() = %v, want %v", gotApiKey, "token")
		}
	})

	t.Run("api key from query", func(t *testing.T) {
		hed := http.Header{}
		url := &url.URL{}
		req := http.Request{Header: hed, URL: url}
		ctx := &gin.Context{Request: &req}
		svc := &cfg.Server{}
		req.URL.RawQuery = "api_key=token"
		gotApiKey, err := ParseAPIKey(ctx, svc)
		if (err != nil) != false {
			t.Errorf("ParseAPIKey() error = %v, wantErr %v", err, false)
			return
		}
		if gotApiKey != "token" {
			t.Errorf("ParseAPIKey() = %v, want %v", gotApiKey, "token")
		}
	})

	t.Run("api key from param", func(t *testing.T) {
		hed := http.Header{}
		url := &url.URL{}
		req := http.Request{Header: hed, URL: url}
		ctx := &gin.Context{Request: &req}
		svc := &cfg.Server{}
		ctx.Params = []gin.Param{
			{
				Key:   "api_key",
				Value: "token",
			},
		}
		gotApiKey, err := ParseAPIKey(ctx, svc)
		if (err != nil) != false {
			t.Errorf("ParseAPIKey() error = %v, wantErr %v", err, false)
			return
		}
		if gotApiKey != "token" {
			t.Errorf("ParseAPIKey() = %v, want %v", gotApiKey, "token")
		}
	})
}

func TestParseToken(t *testing.T) {
	jwtParse = func(tokenString string, keyFunc jwt.Keyfunc) (*jwt.Token, error) {
		return &jwt.Token{Raw: "token"}, nil
	}
	t.Run("jwt token from cookie", func(t *testing.T) {
		hed := http.Header{}
		url := &url.URL{}
		req := http.Request{Header: hed, URL: url}
		ctx := &gin.Context{Request: &req}
		svc := &cfg.Server{
			JWT: cfg.JWT{
				Algorithm: "HS512",
				Secret:    "{JWTsecret}",
			},
		}
		cookie := &http.Cookie{
			Name:   "jwt",
			Value:  "token",
			MaxAge: 300,
		}
		req.AddCookie(cookie)
		gotApiKey, err := ParseToken(ctx, svc)
		if (err != nil) != false {
			t.Errorf("ParseToken() error = %v, wantErr %v", err, false)
			return
		}
		if gotApiKey.Raw != "token" {
			t.Errorf("ParseToken() = %v, want %v", gotApiKey, "token")
		}
	})
	t.Run("token from query", func(t *testing.T) {
		hed := http.Header{}
		url := &url.URL{}
		req := http.Request{Header: hed, URL: url}
		ctx := &gin.Context{Request: &req}
		svc := &cfg.Server{
			JWT: cfg.JWT{
				Algorithm: "HS512",
				Secret:    "{JWTsecret}",
			},
		}
		req.URL.RawQuery = "token=token"
		gotApiKey, err := ParseToken(ctx, svc)
		if (err != nil) != false {
			t.Errorf("ParseToken() error = %v, wantErr %v", err, false)
			return
		}
		if gotApiKey.Raw != "token" {
			t.Errorf("ParseToken() = %v, want %v", gotApiKey, "token")
		}
	})
	t.Run("token from header", func(t *testing.T) {
		hed := http.Header{}
		url := &url.URL{}
		req := http.Request{Header: hed, URL: url}
		ctx := &gin.Context{Request: &req}
		svc := &cfg.Server{
			JWT: cfg.JWT{
				Algorithm: "HS512",
				Secret:    "{JWTsecret}",
			},
		}
		req.Header.Set("Authorization", "Bearer token")
		gotApiKey, err := ParseToken(ctx, svc)
		if (err != nil) != false {
			t.Errorf("ParseToken() error = %v, wantErr %v", err, false)
			return
		}
		if gotApiKey.Raw != "token" {
			t.Errorf("ParseToken() = %v, want %v", gotApiKey, "token")
		}
	})
}

func Test_addToContext(t *testing.T) {
	type testObj struct {
		key   string
		value string
	}
	hed := http.Header{}
	url := &url.URL{}
	req := &http.Request{Header: hed, URL: url}
	ctx := &gin.Context{Request: req}
	type args struct {
		c     *gin.Context
		key   consts.ContextKey
		value interface{}
	}

	tObj := testObj{
		key:   "user",
		value: "id",
	}
	tests := []struct {
		name string
		args args
		want testObj
	}{
		{
			name: "setting context",
			args: args{
				c:     ctx,
				key:   consts.ProjectContextKeys.UserCtxKey,
				value: tObj,
			},
			want: tObj,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addToContext(tt.args.c, tt.args.key, tt.args.value)
			ctxObj, exist := tt.args.c.Get(string(consts.ProjectContextKeys.UserCtxKey))
			if exist != true {
				t.Errorf("addToContext() adds to context got %v, want %v", false, true)
			}
			if !reflect.DeepEqual(ctxObj, tt.want) {
				t.Errorf("addToContext() adds context got %v, want %v", ctxObj, tt.want)
			}
		})
	}
}

func Test_addUserIdToContext(t *testing.T) {
	hed := http.Header{}
	url := &url.URL{}
	req := &http.Request{Header: hed, URL: url}
	ctx := &gin.Context{Request: req}

	uid, _ := uuid.NewV4()
	type args struct {
		c      *gin.Context
		userID uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want uuid.UUID
	}{
		{
			name: "setting context",
			args: args{
				c:      ctx,
				userID: uid,
			},
			want: uid,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addUserIdToContext(tt.args.c, tt.args.userID)
			ctxObj, exist := tt.args.c.Get(string(consts.ProjectContextKeys.UserIDCtxKey))
			if exist != true {
				t.Errorf("addUserIdToContext() adds to context got %v, want %v", false, true)
			}
			if !reflect.DeepEqual(ctxObj, tt.want) {
				t.Errorf("addUserIdToContext() adds context got %v, want %v", ctxObj, tt.want)
			}
		})
	}
}
