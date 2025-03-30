package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/auth"
	"repertoire/server/internal/wrapper"
)

type AuthService interface {
	Refresh(request requests.RefreshRequest) (string, *wrapper.ErrorCode)
	SignIn(request requests.SignInRequest) (string, *wrapper.ErrorCode)
	SignUp(request requests.SignUpRequest) (string, *wrapper.ErrorCode)
	GetCentrifugoToken(token string) (string, *wrapper.ErrorCode)
}

type authService struct {
	refresh            auth.Refresh
	signIn             auth.SignIn
	signUp             auth.SignUp
	getCentrifugoToken auth.GetCentrifugoToken
}

func NewAuthService(
	refresh auth.Refresh,
	signIn auth.SignIn,
	signUp auth.SignUp,
	getCentrifugoToken auth.GetCentrifugoToken,
) AuthService {
	return &authService{
		refresh:            refresh,
		signIn:             signIn,
		signUp:             signUp,
		getCentrifugoToken: getCentrifugoToken,
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

func (a *authService) GetCentrifugoToken(token string) (string, *wrapper.ErrorCode) {
	return a.getCentrifugoToken.Handle(token)
}
