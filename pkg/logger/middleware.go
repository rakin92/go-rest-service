package logger

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rakin92/travel/pkg/consts"
)

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

func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

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
