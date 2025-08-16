package exceptions

import (
	"errors"

	"scheduly.io/core/internal/domain"
)

var (
	ErrAuthInvalidCredentials domain.Error = errors.New("invalid credentials")
	ErrAuthHeaderMissing      domain.Error = errors.New("authentication required")
	ErrInvalidAuthFormat      domain.Error = errors.New("authorization header format must be bearer {token}")
	ErrInvalidToken           domain.Error = errors.New("invalid or expired token")
	ErrSourceExpiredOrInvalid domain.Error = errors.New("session expired or invalid")
	ErrSourceNotFound         domain.Error = errors.New("source not found")
)
