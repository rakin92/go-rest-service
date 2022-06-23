// Package logger is an wrapper around the zerolog library
package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var logger *Logger

func init() {
	logger = NewLogger()
}

// Logger enforces specific log message formats
type Logger struct {
	*zerolog.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var baseLogger = zerolog.New(os.Stdout).With().Timestamp().Stack().Logger()

	return &Logger{&baseLogger}
}

// Event stores messages to log later, from our standard interface
type Event struct {
	id      int
	message string
}

// Declare variables to store log messages as new Events
var (
	invalidArgMessage      = Event{1, "Invalid arg: %s"}
	invalidArgValueMessage = Event{2, "Invalid value for argument: %s: %v"}
	missingArgMessage      = Event{3, "Missing arg: %s"}
)

// Expose some log functions:

// Errorfn Log errors of a [fn] with format
func Errorfn(fn string, err error) error {
	outerr := fmt.Errorf("[%s]: %v", fn, err)
	logger.Error().Err(outerr).Msg(outerr.Error())
	return outerr
}

// InvalidArg is a standard error message for invalid argument
func InvalidArg(argumentName string) error {
	outerr := fmt.Errorf(invalidArgMessage.message, argumentName)
	logger.Error().Err(outerr).Msg(outerr.Error())
	return outerr
}

// InvalidArgValue is a standard error message for missing argument value
func InvalidArgValue(argumentName string, argumentValue string) error {
	outerr := fmt.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
	logger.Error().Err(outerr).Msg(outerr.Error())
	return outerr
}

// MissingArg is a standard error message for missing arguments
func MissingArg(argumentName string) error {
	outerr := fmt.Errorf(missingArgMessage.message, argumentName)
	log.Error().Err(outerr).Msg(outerr.Error())
	return outerr
}

// Debug logs a new message with debug level.
func Debug(message string, args ...any) {
	logger.Debug().Msgf(message, args...)
}

// Info logs a new message with info level.
func Info(message string, args ...any) {
	logger.Info().Msgf(message, args...)
}

// Warn logs a new message with warn level.
func Warn(message string, args ...any) {
	logger.Warn().Msgf(message, args...)
}

// Error logs a new message with error level.
func Error(err *error, message string, args ...any) {
	if err != nil {
		logger.Error().Err(*err).Msgf(message, args...)
	} else {
		logger.Error().Msgf(message, args...)
	}
}

// Fatal logs a new message with fatal level. The os.Exit(1) function
// is called which terminates the program immediately.
func Fatal(err *error, message string, args ...any) {
	if err != nil {
		logger.Fatal().Err(*err).Msgf(message, args...)
	} else {
		logger.Fatal().Msgf(message, args...)
	}
}

// Panic logs a new message with panic level. The panic() function
// is called by the Msg method, which stops the ordinary flow of a goroutine.
func Panic(err *error, message string, args ...any) {
	if err != nil {
		logger.Panic().Err(*err).Msgf(message, args...)
	} else {
		logger.Panic().Msgf(message, args...)
	}
}
