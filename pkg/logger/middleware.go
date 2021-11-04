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

// ErrorLogger is a handler function for any gin error type
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT is a handler function for any gin error type
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.Writer.Written() {
			json := c.Errors.ByType(typ).JSON()
			if json != nil {
				c.JSON(-1, json)
			}
		}
	}
}

// Middleware to log our gin requests in formatted JSON.
// Can be added to be used by our Gin router.
func Middleware(serName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		// before request
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()

		if raw != "" {
			path = path + "?" + raw
		}

		lf := &logFields{
			SerName:    serName,
			Path:       path,
			Latency:    time.Since(t),
			Method:     c.Request.Method,
			StatusCode: c.Writer.Status(),
			ClientIP:   c.ClientIP(),
			MsgStr:     c.Errors.String(),
		}

		lf.User = "anonymous"
		u, exist := c.Get(string(consts.ProjectContextKeys.UserCtxKey))
		if exist {
			lf.User = fmt.Sprintf("%v", u)
		}

		logSwitch(lf)
	}
}

// logSwitch logs different levels of logs based on status code
func logSwitch(lf *logFields) {
	// TODO - look into using log with fields
	switch {
	case lf.StatusCode >= 400 && lf.StatusCode < 500:
		{
			logger.Warn().Str("user", lf.User).Str("service_name", lf.SerName).Str("method", lf.Method).Str("path", lf.Path).Dur("latency", lf.Latency).Int("status", lf.StatusCode).Str("client_ip", lf.ClientIP).Msg(lf.MsgStr)
		}
	case lf.StatusCode >= 500:
		{
			logger.Error().Str("user", lf.User).Str("service_name", lf.SerName).Str("method", lf.Method).Str("path", lf.Path).Dur("latency", lf.Latency).Int("status", lf.StatusCode).Str("client_ip", lf.ClientIP).Msg(lf.MsgStr)
		}
	default:
		logger.Info().Str("user", lf.User).Str("service_name", lf.SerName).Str("method", lf.Method).Str("path", lf.Path).Dur("latency", lf.Latency).Int("status", lf.StatusCode).Str("client_ip", lf.ClientIP).Msg(lf.MsgStr)
	}
}
