package local

import (
	"context"
	"log"

	"buke.io/core/internal/domain"
	"buke.io/core/internal/utils/logger"
)

type Logger struct {
	logger *logger.Logger
}

func NewLocalLogger() *Logger {
	return &Logger{
		logger: logger.New("debug"),
	}
}

type LogPayload struct {
	Message string                 `json:"message,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (instance *Logger) Debug(ctx context.Context, v ...interface{}) {

	log.Print(v...)
}

func (instance *Logger) DebugWithVar(ctx context.Context, message string, vars map[string]interface{}) {
	log.Print(message)
}

func (instance *Logger) Info(ctx context.Context, v ...interface{}) {
	log.Print(v...)
}

func (instance *Logger) Warning(ctx context.Context, v ...interface{}) {
	log.Print(v...)
}

func (instance *Logger) ErrorWithVar(ctx context.Context, err error, vars map[string]interface{}) {
	log.Print(err)
	//debug.PrintStack()
}

func (instance *Logger) Error(ctx context.Context, v ...interface{}) {
	log.Print(v...)
	//debug.PrintStack()
}

func (instance *Logger) LogHttpReq(data *domain.LogHttpEntry) {
	log.Print(data)
}
