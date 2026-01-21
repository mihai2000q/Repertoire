package service

import (
	"errors"
	"net/http"
	"reflect"
	"repertoire/auth/api/requests"
	"repertoire/auth/data/logger"
	"repertoire/auth/data/repository"
	"repertoire/auth/data/service"
	"repertoire/auth/internal/wrapper"
	"repertoire/auth/model"
	"strings"

	"go.uber.org/zap"
)

type MainService interface {
	Refresh(request requests.RefreshRequest) (string, *wrapper.ErrorCode)
	SignIn(request requests.SignInRequest) (string, *wrapper.ErrorCode)
}

type mainService struct {
	jwtService     service.JwtService
	bCryptService  service.BCryptService
	userRepository repository.UserRepository
	logger         *logger.Logger
}

func NewMainService(
	jwtService service.JwtService,
	bCryptService service.BCryptService,
	userRepository repository.UserRepository,
	logger *logger.Logger,
) MainService {
	return &mainService{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
		logger:         logger,
	}
}

func (m *mainService) Refresh(request requests.RefreshRequest) (string, *wrapper.ErrorCode) {
	// validate token
	userID, errCode := m.jwtService.Validate(request.Token)
	if errCode != nil {
		if errCode.Code == http.StatusUnauthorized {
			m.logger.Warn("Invalid token", zap.Error(errCode.Error), zap.String("token", request.Token))
			return "", wrapper.UnauthorizedError(errors.New("invalid token"))
		}
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
