package exceptions

import (
	"errors"

	"buke.io/core/internal/domain"
)

var (
	ErrUserNotFound domain.Error = errors.New("user not found")
)
