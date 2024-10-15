package service

import (
	"repertoire/api/requests"
	"repertoire/domain/usecases/auth"
	"repertoire/utils"
)

type AuthService interface {
	Refresh(request requests.RefreshRequest) (string, *utils.ErrorCode)
	SignIn(request requests.SignInRequest) (string, *utils.ErrorCode)
	SignUp(request requests.SignUpRequest) (string, *utils.ErrorCode)
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

func (a *authService) Refresh(request requests.RefreshRequest) (string, *utils.ErrorCode) {
	return a.refresh.Handle(request)
}

func (a *authService) SignIn(request requests.SignInRequest) (string, *utils.ErrorCode) {
	return a.signIn.Handle(request)
}

func (a *authService) SignUp(request requests.SignUpRequest) (string, *utils.ErrorCode) {
	return a.signUp.Handle(request)
}
