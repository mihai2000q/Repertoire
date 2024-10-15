package auth

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
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
) *SignIn {
	return &SignIn{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
}

func (s *SignIn) Handle(request requests.SignInRequest) (string, *utils.ErrorCode) {
	// get user
	var user models.User
	err := s.userRepository.GetByEmail(&user, strings.ToLower(request.Email))
	if err != nil {
		return "", utils.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return "", utils.UnauthorizedError(errors.New("invalid credentials"))
	}

	// check password
	err = s.bCryptService.CompareHash(user.Password, request.Password)
	if err != nil {
		return "", utils.UnauthorizedError(errors.New("invalid credentials"))
	}

	return s.jwtService.CreateToken(user)
}
