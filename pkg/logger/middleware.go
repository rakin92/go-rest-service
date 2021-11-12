package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakin92/go-rest-service/pkg/consts"
)

// logFields is used to structure of our request log data
type logFields struct {
	SerName    string
	Path       string
	Latency    time.Duration
	Method     string
	StatusCode int
	ClientIP   string
	MsgStr     string
	User       string
}

func prepareLogFields(c *gin.Context, path, rawQ, server string, t time.Duration) *logFields {
	if rawQ != "" {
		path = path + "?" + rawQ
	}

	lf := &logFields{
		SerName:    server,
		Path:       path,
		Latency:    t,
		Method:     c.Request.Method,
		StatusCode: c.Writer.Status(),
		ClientIP:   c.ClientIP(),
		MsgStr:     c.Errors.String(),
	}

	lf.User = "anonymous"
	u, exist := c.Get(string(consts.ProjectContextKeys.UserIDCtxKey))
	if exist {
		lf.User = fmt.Sprintf("%v", u)
	}
	return lf
}

// Middleware to log our gin requests in formatted JSON.
// Can be added to be used by our Gin router.
func Middleware(serverName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		path := c.Request.URL.Path
		rawQ := c.Request.URL.RawQuery
		c.Next()

		d := time.Since(t)
		fields := prepareLogFields(c, path, rawQ, serverName, d)
		logSwitch(fields)
	}
}

// logSwitch logs different levels of logs based on status code
func logSwitch(lf *logFields) error {
	// TODO: look into using log with fields
	switch {
	case lf.StatusCode >= 400 && lf.StatusCode < 500:
		{
			logger.Warn().Str("user", lf.User).Str("service_name", lf.SerName).Str("method", lf.Method).Str("path", lf.Path).Dur("latency", lf.Latency).Int("status", lf.StatusCode).Str("client_ip", lf.ClientIP).Msg(lf.MsgStr)
			return fmt.Errorf("received status code of 400s: %d", lf.StatusCode)
		}
	case lf.StatusCode >= 500:
		{
			logger.Error().Str("user", lf.User).Str("service_name", lf.SerName).Str("method", lf.Method).Str("path", lf.Path).Dur("latency", lf.Latency).Int("status", lf.StatusCode).Str("client_ip", lf.ClientIP).Msg(lf.MsgStr)
			return fmt.Errorf("received status code of 500s: %d", lf.StatusCode)
		}
	default:
		logger.Info().Str("user", lf.User).Str("service_name", lf.SerName).Str("method", lf.Method).Str("path", lf.Path).Dur("latency", lf.Latency).Int("status", lf.StatusCode).Str("client_ip", lf.ClientIP).Msg(lf.MsgStr)
	}
	return nil
}
