package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// SetGlobalLevel задать глобальный уровень логирования
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

// Logger zerolog логер
type Logger struct {
	logger *zerolog.Logger
}

// Debug записать сообщение с уровнем debug
func (l *Logger) Debug(message string, args ...interface{}) {
	l.logger.Debug().Msgf(message, args...)
}

// Info записать сообщение с уровнем info
func (l *Logger) Info(message string, args ...interface{}) {
	l.logger.Info().Msgf(message, args...)
}

// Warn записать сообщение с уровнем warn
func (l *Logger) Warn(message string, args ...interface{}) {
	l.logger.Warn().Msgf(message, args...)
}

// Error записать сообщение с уровнем error
func (l *Logger) Error(err error) {
	l.logger.Error().Stack().Err(err).Msg("")
}

// Fatal записать сообщение с уровнем fatal и завершить работу
func (l *Logger) Fatal(err error) {
	l.logger.Error().Stack().Err(err).Msg("")
	os.Exit(1)
}

// LoggerOfComponent получить логер для компонента
func LoggerOfComponent(component string) Interface {
	logger := log.With().Str("component", component).Logger()
	return &Logger{logger: &logger}

}
