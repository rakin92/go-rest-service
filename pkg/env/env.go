// Package env holds all environment variable related code
package env

import (
	"os"
	"strconv"

	"github.com/rakin92/go-rest-service/pkg/logger"
)

// MustGet will return the env or panic if it is not present
func MustGet(k string) string {
	v := os.Getenv(k)
	if v == "" {
		logger.MissingArg(k)
		logger.Panicf(nil, "ENV missing, key: "+k)
	}
	return v
}

// MustGetBool will return the env as boolean or panic if it is not present
func MustGetBool(k string) bool {
	v := os.Getenv(k)
	if v == "" {
		logger.MissingArg(k)
		logger.Panicf(nil, "ENV missing, key: %s", k)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		logger.MissingArg(k)
		logger.Panicf(&err, "ENV err: [%s]", err.Error())
	}
	return b
}

// MustGetInt32 will return the env as integer or panic if it is not present
func MustGetInt32(k string) int {
	v := os.Getenv(k)
	if v == "" {
		logger.MissingArg(k)
		logger.Panicf(nil, "ENV missing, key: %s", k)
	}
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		logger.MissingArg(k)
		logger.Panicf(&err, "ENV err: [%s]", err.Error())
	}
	return int(i)
}

// MustGetInt64 will return the env as int64 or panic if it is not present
func MustGetInt64(k string) int64 {
	v := os.Getenv(k)
	if v == "" {
		logger.MissingArg(k)
		logger.Panicf(nil, "ENV missing, key: %s", k)
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logger.MissingArg(k)
		logger.Panicf(&err, "ENV err: [%s]"+err.Error())
	}
	return i
}
