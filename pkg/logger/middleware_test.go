package logger

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakin92/go-rest-service/pkg/consts"
)

func Test_prepareLogFields(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Method: "POST",
	}
	ctx.Request.Response = &http.Response{
		StatusCode: 200,
	}

	now := time.Now()
	d := time.Since(now)

	type args struct {
		c      *gin.Context
		rawQ   string
		path   string
		server string
		t      time.Duration
		user   string
	}
	tests := []struct {
		name string
		args args
		want *logFields
	}{
		{
			name: "passing with anonymous",
			args: args{
				c:      ctx,
				path:   "/user",
				rawQ:   "",
				server: "test",
				t:      d,
				user:   "",
			},
			want: &logFields{
				SerName:    "test",
				Path:       "/user",
				Latency:    d,
				Method:     "POST",
				StatusCode: 200,
				ClientIP:   "",
				MsgStr:     "",
				User:       "anonymous",
			},
		},
		{
			name: "passing with query params",
			args: args{
				c:      ctx,
				path:   "/user",
				rawQ:   "id=1",
				server: "test",
				t:      d,
				user:   "",
			},
			want: &logFields{
				SerName:    "test",
				Path:       "/user?id=1",
				Latency:    d,
				Method:     "POST",
				StatusCode: 200,
				ClientIP:   "",
				MsgStr:     "",
				User:       "anonymous",
			},
		},
		{
			name: "passing with user_id",
			args: args{
				c:      ctx,
				path:   "/user",
				rawQ:   "",
				server: "test",
				t:      d,
				user:   "user_id",
			},
			want: &logFields{
				SerName:    "test",
				Path:       "/user",
				Latency:    d,
				Method:     "POST",
				StatusCode: 200,
				ClientIP:   "",
				MsgStr:     "",
				User:       "user_id",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.user != "" {
				tt.args.c.Set(string(consts.ProjectContextKeys.UserIDCtxKey), tt.args.user)
			}
			if got := prepareLogFields(tt.args.c, tt.args.path, tt.args.rawQ, tt.args.server, tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prepareLogFields() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_logSwitch(t *testing.T) {
	type args struct {
		lf *logFields
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "400s",
			args: args{
				lf: &logFields{
					StatusCode: 404,
				},
			},
			wantErr: true,
		},
		{
			name: "500s",
			args: args{
				lf: &logFields{
					StatusCode: 503,
				},
			},
			wantErr: true,
		},
		{
			name: "200s",
			args: args{
				lf: &logFields{
					StatusCode: 200,
				},
			},
			wantErr: false,
		},
		{
			name: "300s",
			args: args{
				lf: &logFields{
					StatusCode: 301,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := logSwitch(tt.args.lf); (err != nil) != tt.wantErr {
				t.Errorf("logSwitch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
