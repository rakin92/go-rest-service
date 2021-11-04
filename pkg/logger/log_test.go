package logger_test

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/rakin92/go-rest-service/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {

	t.Run("Returns a valid *Logger", func(t *testing.T) {
		got := logger.NewLogger()
		assert.IsType(t, reflect.TypeOf(got), reflect.TypeOf(logger.Logger{}))
	})
}

func TestPanic(t *testing.T) {
	e := errors.New("an error")
	type args struct {
		err     *error
		message string
		args    []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "panic with error",
			args: args{
				err:     &e,
				message: "panic",
			},
		},
		{
			name: "panic without error",
			args: args{
				err:     nil,
				message: "panic",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			func() {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Panic should have panicked!")
					}
				}()
				logger.Panic(tt.args.err, tt.args.message, tt.args.args...)
			}()
		})
	}
}

func TestFatal(t *testing.T) {
	e := errors.New("an error")
	type args struct {
		err     *error
		message string
		args    []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "fatal with error",
			args: args{
				err:     &e,
				message: "fatal",
			},
		},
		{
			name: "fatal without error",
			args: args{
				err:     nil,
				message: "fatal",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if os.Getenv("FLAG") == "1" {
				logger.Fatal(tt.args.err, tt.args.message, tt.args.args...)
				return
			}
			// Run the test in a subprocess
			cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
			cmd.Env = append(os.Environ(), "FLAG=1")
			err := cmd.Run()

			// Cast the error as *exec.ExitError and compare the result
			e, ok := err.(*exec.ExitError)
			expectedErrorString := "exit status 1"
			assert.Equal(t, true, ok)
			assert.Equal(t, expectedErrorString, e.Error())
		})
	}
}

func TestError(t *testing.T) {
	e := errors.New("an error")
	type args struct {
		err     *error
		message string
		args    []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "error with error",
			args: args{
				err:     &e,
				message: "error",
			},
		},
		{
			name: "error without error",
			args: args{
				err:     nil,
				message: "error",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.Error(tt.args.err, tt.args.message, tt.args.args...)
		})
	}
}

func TestMissingArg(t *testing.T) {
	e := errors.New("error")
	got := logger.MissingArg("type")
	assert.IsType(t, reflect.TypeOf(got), reflect.TypeOf(e))
	assert.Equal(t, got.Error(), "Missing arg: type")
}

func TestInvalidArg(t *testing.T) {
	e := errors.New("error")
	got := logger.InvalidArg("type")
	assert.IsType(t, reflect.TypeOf(got), reflect.TypeOf(e))
	assert.Equal(t, got.Error(), "Invalid arg: type")
}

func TestInvalidArgValue(t *testing.T) {
	e := errors.New("error")
	got := logger.InvalidArgValue("type", "error")
	assert.IsType(t, reflect.TypeOf(got), reflect.TypeOf(e))
	assert.Equal(t, got.Error(), "Invalid value for argument: type: error")
}
