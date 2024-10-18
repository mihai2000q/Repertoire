package auth

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/request"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
	"strings"
)

type SignIn struct {
	jwtService     service.JwtService
	bCryptService  service.BCryptService
	userRepository repository.UserRepository
}

func NewSignIn(
	jwtService service.JwtService,
	bCryptService service.BCryptService,
	userRepository repository.UserRepository,
) SignIn {
	return SignIn{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
}

func (s *SignIn) Handle(request request.SignInRequest) (string, *wrapper.ErrorCode) {
	// get user
	var user model.User
	err := s.userRepository.GetByEmail(&user, strings.ToLower(request.Email))
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return "", wrapper.UnauthorizedError(errors.New("invalid credentials"))
	}

	// check password
	err = s.bCryptService.CompareHash(user.Password, request.Password)
	if err != nil {
		return "", wrapper.UnauthorizedError(errors.New("invalid credentials"))
	}

	return s.jwtService.CreateToken(user)
}
