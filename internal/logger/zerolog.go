package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetGlobalLevel(level string) {
	var lvl zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		lvl = zerolog.ErrorLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "debug":
		lvl = zerolog.DebugLevel
	default:
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)
}

type Logger struct {
	logger *zerolog.Logger
}

func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug().Msgf(message, args...)
}

func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info().Msgf(message, args...)
}

func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn().Msgf(message, args...)
}

func (l *Logger) Error(err error) {
	l.logger.Error().Stack().Err(err).Msg("")
}

func (l *Logger) Fatal(err error) {
	l.logger.Error().Stack().Err(err).Msg("")
	os.Exit(1)
}

func LoggerOfComponent(component string) Interface {
	logger := log.With().Str("component", component).Logger()
	return &Logger{logger: &logger}

}
