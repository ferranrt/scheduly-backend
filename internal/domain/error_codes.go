package domain

import "errors"

var (
	ErrEnvInvalid Error = errors.New("ENV_INVALID")

	ErrBadRequestPayload Error = errors.New("BAD_REQUEST_PAYLOAD")

	ErrInternalServerError Error = errors.New("SERVER_ERROR")
	ErrForbidden           Error = errors.New("FORBIDDEN")
	ErrInsert              Error = errors.New("INSERT_ERROR")
	ErrUpdate              Error = errors.New("UPDATE_ERROR")
	ErrDelete              Error = errors.New("DELETE_ERROR")
	ErrNotFound            Error = errors.New("NOT_FOUND")
	ErrBadRequest          Error = errors.New("BAD_REQUEST")
	ErrUnexpectedError     Error = errors.New("UNEXPECTED_ERROR")
)
