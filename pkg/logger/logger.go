package logger

import (
	"log"
	"os"
)

// Logger methods interface
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatal(v ...interface{})
}

type apiLogger struct {
	logger *log.Logger
}

func (a *apiLogger) Println(v ...interface{}) {
	a.logger.Println(v...)
}

func (a *apiLogger) Printf(format string, v ...interface{}) {
	a.logger.Printf(format, v...)
}

func (a *apiLogger) Fatalf(format string, v ...interface{}) {
	a.logger.Fatalf(format, v...)
}

func (a *apiLogger) Fatal(v ...interface{}) {
	a.logger.Fatal(v...)
}

func NewApiLogger() Logger {
	return &apiLogger{logger: log.New(os.Stdout, "TODOLOG", log.Llongfile)}
}