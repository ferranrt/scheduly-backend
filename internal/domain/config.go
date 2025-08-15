package domain

import "time"

type JWTConfig struct {
	AtkSecret string
	RtkSecret string
	Expiry    time.Duration
}
