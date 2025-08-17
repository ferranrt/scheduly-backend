package mocks

import (
	"context"
	"log"

	"buke.io/core/internal/domain"
	"github.com/stretchr/testify/mock"
)

type LoggerMock struct {
	mock.Mock
}

func (p *LoggerMock) Debug(ctx context.Context, v ...interface{}) {
	log.Print(v...)
}

func (p *LoggerMock) DebugWithVar(ctx context.Context, message string, vars map[string]interface{}) {
	log.Print(message)
}

func (p *LoggerMock) Info(ctx context.Context, v ...interface{}) {
	log.Print(v...)
}

func (p *LoggerMock) Warning(ctx context.Context, v ...interface{}) {
	log.Print(v...)
}

func (p *LoggerMock) ErrorWithVar(ctx context.Context, err error, vars map[string]interface{}) {

}

func (p *LoggerMock) Error(ctx context.Context, v ...interface{}) {

}

func (l *LoggerMock) LogHttpReq(data *domain.LogHttpEntry) {
	log.Print(data)
}
