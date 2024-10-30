package auth

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
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

func (s *SignIn) Handle(request requests.SignInRequest) (string, *wrapper.ErrorCode) {
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
