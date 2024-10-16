package applog

import "log"

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
	StandardLogger(level Level) *log.Logger
}
