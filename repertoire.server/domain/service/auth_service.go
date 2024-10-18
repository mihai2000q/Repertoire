package service

import (
	"repertoire/api/request"
	"repertoire/domain/usecase/auth"
	"repertoire/utils/wrapper"
)

type AuthService interface {
	Refresh(request request.RefreshRequest) (string, *wrapper.ErrorCode)
	SignIn(request request.SignInRequest) (string, *wrapper.ErrorCode)
	SignUp(request request.SignUpRequest) (string, *wrapper.ErrorCode)
}

type authService struct {
	refresh auth.Refresh
	signIn  auth.SignIn
	signUp  auth.SignUp
}

func NewAuthService(
	refresh auth.Refresh,
	signIn auth.SignIn,
	signUp auth.SignUp,
) AuthService {
	return &authService{
		refresh: refresh,
		signIn:  signIn,
		signUp:  signUp,
	}
}

func (a *authService) Refresh(request request.RefreshRequest) (string, *wrapper.ErrorCode) {
	return a.refresh.Handle(request)
}

func (a *authService) SignIn(request request.SignInRequest) (string, *wrapper.ErrorCode) {
	return a.signIn.Handle(request)
}

func (a *authService) SignUp(request request.SignUpRequest) (string, *wrapper.ErrorCode) {
	return a.signUp.Handle(request)
}
