package exceptions

import (
	"errors"

	"scheduly.io/core/internal/domain"
)

var (
	ErrAuthInvalidCredentials     domain.Error = errors.New("invalid credentials")
	ErrAuthHeaderMissing          domain.Error = errors.New("authentication required")
	ErrAuthInvalidAuthFormat      domain.Error = errors.New("authorization header format must be bearer {token}")
	ErrAuthInvalidToken           domain.Error = errors.New("invalid or expired token")
	ErrAuthSourceExpiredOrInvalid domain.Error = errors.New("session expired or invalid")
	ErrAuthSourceNotFound         domain.Error = errors.New("source not found")
)
