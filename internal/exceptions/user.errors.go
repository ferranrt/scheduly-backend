package exceptions

import (
	"errors"

	"bifur.app/core/internal/domain"
)

var (
	ErrUserNotFound domain.Error = errors.New("user not found")
)
