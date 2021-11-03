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

// StandardLogger enforces specific log message formats
type Logger struct {
	*zerolog.Logger
}

// NewLogger initializes the standard logger
func NewLogger() *Logger {
	var baseLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

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

// Debug Log
func Debug(message string) {
	logger.Debug().Msg(message)
}

// Debugf Log
func Debugf(message string, args ...interface{}) {
	logger.Debug().Msgf(message, args...)
}

// Errorfn Log errors of a [fn] with format
func Errorfn(fn string, err error) error {
	outerr := fmt.Errorf("[%s]: %v", fn, err)
	logger.Error().Stack().Err(outerr).Msg(outerr.Error())
	return outerr
}

// InvalidArg is a standard error message
func InvalidArg(argumentName string) error {
	outerr := fmt.Errorf(invalidArgMessage.message, argumentName)
	logger.Error().Stack().Err(outerr).Msg(outerr.Error())
	return outerr
}

// InvalidArgValue is a standard error message
func InvalidArgValue(argumentName string, argumentValue string) error {
	outerr := fmt.Errorf(invalidArgValueMessage.message, argumentName, argumentValue)
	logger.Error().Stack().Err(outerr).Msg(outerr.Error())
	return outerr
}

// MissingArg is a standard error message
func MissingArg(argumentName string) error {
	outerr := fmt.Errorf(missingArgMessage.message, argumentName)
	log.Error().Stack().Err(outerr).Msg(outerr.Error())
	return outerr
}

// Info Log
func Info(message string) {
	logger.Info().Msg(message)
}

// Infof Log
func Infof(message string, args ...interface{}) {
	logger.Info().Msgf(message, args...)
}

// Warn Log
func Warn(message string) {
	logger.Warn().Msg(message)
}

// Warnf Log
func Warnf(message string, args ...interface{}) {
	logger.Warn().Msgf(message, args...)
}

// Panic Log
func Panic(err *error, message string) {
	if err != nil {
		logger.Panic().Stack().Err(*err).Msg(message)
	} else {
		logger.Panic().Stack().Msg(message)
	}
}

// Panicf Log
func Panicf(err *error, message string, args ...interface{}) {
	if err != nil {
		logger.Panic().Stack().Err(*err).Msgf(message, args)
	} else {
		logger.Panic().Stack().Msgf(message, args)
	}
}

// Error Log
func Error(err *error, message string) {
	if err != nil {
		logger.Error().Stack().Err(*err).Msg(message)
	} else {
		logger.Error().Stack().Msg(message)
	}
}

// Errorf Log
func Errorf(err *error, message string, args ...interface{}) {
	if err != nil {
		logger.Error().Stack().Err(*err).Msgf(message, args)
	} else {
		logger.Error().Stack().Msgf(message, args)
	}
}

// Fatal Log
func Fatal(err *error, message string) {
	if err != nil {
		logger.Fatal().Stack().Err(*err).Msg(message)
	} else {
		logger.Fatal().Stack().Msg(message)
	}
}

// Fatalf Log
func Fatalf(err *error, message string, args ...interface{}) {
	if err != nil {
		logger.Fatal().Stack().Err(*err).Msgf(message, args)
	} else {
		logger.Fatal().Stack().Msgf(message, args)
	}
}
