package ports

import (
	"context"

	"buke.io/core/internal/domain"
)

type Logger interface {
	Debug(ctx context.Context, v ...interface{})
	DebugWithVar(ctx context.Context, message string, vars map[string]interface{})
	Info(ctx context.Context, v ...interface{})
	Warning(ctx context.Context, v ...interface{})
	Error(ctx context.Context, v ...interface{})
	ErrorWithVar(ctx context.Context, err error, vars map[string]interface{})

	LogHttpReq(data *domain.LogHttpEntry)
}
