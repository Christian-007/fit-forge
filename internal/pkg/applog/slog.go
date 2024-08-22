package applog

import (
	"log/slog"
)

type SlogLogger struct {
	logger *slog.Logger
}

func New(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{
		logger: logger,
	}
}

func (l *SlogLogger) Error(msg string, ctx ...interface{}) {
	l.logger.Error(msg, ctx...)
}
