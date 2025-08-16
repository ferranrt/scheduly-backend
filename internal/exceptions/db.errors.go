package exceptions

import (
	"errors"

	"scheduly.io/core/internal/domain"
)

var (
	ErrDbInsert domain.Error = errors.New("DB_INSERT_ERROR")
	ErrDbUpdate domain.Error = errors.New("DB_UPDATE_ERROR")
	ErrDbSelect domain.Error = errors.New("DB_SELECT_ERROR")
	ErrDbDelete domain.Error = errors.New("DB_DELETE_ERROR")
)
