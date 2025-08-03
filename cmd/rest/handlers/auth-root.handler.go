package handlers

import "scheduly.io/core/internal/ports/usecases"

type AuthRootHandler struct {
	authUseCase usecases.AuthUseCase
}

func NewAuthRootHandler(authUseCase usecases.AuthUseCase) *AuthRootHandler {
	return &AuthRootHandler{
		authUseCase: authUseCase,
	}
}
