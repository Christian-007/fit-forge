package applog

//go:generate mockgen -source=interface.go -destination=mocks/interface.go
type Logger interface {
	Error(msg string, ctx ...interface{})
	Info(msg string, args ...any)
}
