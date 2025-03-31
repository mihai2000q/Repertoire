package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/auth"
	"repertoire/server/internal/wrapper"
)

type AuthService interface {
	SignUp(request requests.SignUpRequest) (string, *wrapper.ErrorCode)
}

type authService struct {
	signUp auth.SignUp
}

func NewAuthService(
	signUp auth.SignUp,
) AuthService {
	return &authService{
		signUp: signUp,
	}
}

func (a *authService) SignUp(request requests.SignUpRequest) (string, *wrapper.ErrorCode) {
	return a.signUp.Handle(request)
}
