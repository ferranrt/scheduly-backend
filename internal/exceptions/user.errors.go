package exceptions

import (
	"errors"

	"scheduly.io/core/internal/domain"
)

var (
	ErrUserNotFound domain.Error = errors.New("user not found")
)
