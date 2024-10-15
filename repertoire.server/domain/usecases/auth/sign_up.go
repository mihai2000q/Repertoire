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

type SignUp struct {
	jwtService     service.JwtService
	bCryptService  service.BCryptService
	userRepository repository.UserRepository
}

func NewSignUp(
	jwtService service.JwtService,
	bCryptService service.BCryptService,
	userRepository repository.UserRepository,
) *SignUp {
	return &SignUp{
		jwtService:     jwtService,
		bCryptService:  bCryptService,
		userRepository: userRepository,
	}
}

func (s *SignUp) Handle(request requests.SignUpRequest) (string, *utils.ErrorCode) {
	var user models.User

	// check if the user already exists
	email := strings.ToLower(request.Email)
	err := s.userRepository.GetByEmail(&user, email)
	if err != nil {
		return "", utils.InternalServerError(err)
	}
	if user.ID != uuid.Nil {
		return "", utils.BadRequestError(errors.New("user already exists"))
	}

	// hash the password
	hashedPassword, err := s.bCryptService.Hash(request.Password)
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	// create user
	user = models.User{
		ID:       uuid.New(),
		Name:     request.Name,
		Email:    email,
		Password: hashedPassword,
	}
	err = s.userRepository.Create(&user)
	if err != nil {
		return "", utils.InternalServerError(err)
	}

	return s.jwtService.CreateToken(user)
}
