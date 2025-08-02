package domain

import (
	"net/http"
)

type Error error
type HTTPErrorBody struct {
	Code   string `json:"code"`
	Errors any    `json:"errors"`
}

type DomainError struct {
	HTTPCode      int
	HTTPErrorBody HTTPErrorBody
	OriginalError string
}

func NewDomainError(httpCode int, errorCode string, errorMsg any, err error) *DomainError {
	return &DomainError{
		HTTPCode:      httpCode,
		OriginalError: err.Error(),
		HTTPErrorBody: HTTPErrorBody{
			Code:   errorCode,
			Errors: errorMsg,
		},
	}
}

func NewInternalError(err error) *DomainError {
	return &DomainError{
		HTTPCode:      http.StatusInternalServerError,
		OriginalError: err.Error(),
		HTTPErrorBody: HTTPErrorBody{
			Code:   ErrUnexpectedError.Error(),
			Errors: "Internal server error",
		},
	}
}
