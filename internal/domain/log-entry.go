package domain

import (
	"net/http"
	"time"
)

type LogHttpReq struct {
	Request      *http.Request
	RequestSize  int64
	ResponseSize int64
	Status       int
	RemoteIP     string
	Latency      time.Duration
}
type LogTraceAuth struct {
	RawBody string            `json:"rawBody"`
	Headers map[string]string `json:"headers"`
}

type LogHttpEntry struct {
	HTTPRequest *LogHttpReq
	Payload     *LogTraceAuth
}
