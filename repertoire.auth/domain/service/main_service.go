package service

import (
	"errors"
	"reflect"
	"repertoire/auth/api/requests"
	"repertoire/auth/data/repository"
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"
	"strings"
)

type MainService interface {
	Refresh(request requests.RefreshRequest) (string, *wrapper.ErrorCode)
	SignIn(request requests.SignInRequest) (string, *wrapper.ErrorCode)
}

type mainService struct {
	jwtService     service.JwtService
	bCryptService  service.BCryptService
	userRepository repository.UserRepository
}

func NewMainService(
	jwtService service.JwtService,
	bCryptService service.BCryptService,
	userRepository repository.UserRepository,
) MainService {
	return &mainService{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
}

func (m *mainService) Refresh(request requests.RefreshRequest) (string, *wrapper.ErrorCode) {
	// validate token
	userID, errCode := m.jwtService.Validate(request.Token)
	if errCode != nil {
		return "", errCode
	}

	// get user
	var user model.User
	err := m.userRepository.Get(&user, userID)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return "", wrapper.UnauthorizedError(errors.New("not authorized"))
	}

	return m.jwtService.CreateToken(user)
}

func (m *mainService) SignIn(request requests.SignInRequest) (string, *wrapper.ErrorCode) {
	// get user
	var user model.User
	err := m.userRepository.GetByEmail(&user, strings.ToLower(request.Email))
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return "", wrapper.UnauthorizedError(errors.New("invalid credentials"))
	}

	// check password
	err = m.bCryptService.CompareHash(user.Password, request.Password)
	if err != nil {
		return "", wrapper.UnauthorizedError(errors.New("invalid credentials"))
	}

	return m.jwtService.CreateToken(user)
}
