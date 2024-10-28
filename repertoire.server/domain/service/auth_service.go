package service

import (
	"repertoire/api/requests"
	"repertoire/domain/usecase/auth"
	"repertoire/utils/wrapper"
)

type AuthService interface {
	Refresh(request requests.RefreshRequest) (string, *wrapper.ErrorCode)
	SignIn(request requests.SignInRequest) (string, *wrapper.ErrorCode)
	SignUp(request requests.SignUpRequest) (string, *wrapper.ErrorCode)
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

func (a *authService) Refresh(request requests.RefreshRequest) (string, *wrapper.ErrorCode) {
	return a.refresh.Handle(request)
}

func (a *authService) SignIn(request requests.SignInRequest) (string, *wrapper.ErrorCode) {
	return a.signIn.Handle(request)
}

func (a *authService) SignUp(request requests.SignUpRequest) (string, *wrapper.ErrorCode) {
	return a.signUp.Handle(request)
}
