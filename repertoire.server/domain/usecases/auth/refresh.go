package auth

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type Refresh struct {
	jwtService     service.JwtService
	userRepository repository.UserRepository
}

func NewRefresh(jwtService service.JwtService, userRepository repository.UserRepository) Refresh {
	return Refresh{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
}

func (r *Refresh) Handle(request requests.RefreshRequest) (string, *wrapper.ErrorCode) {
	// validate token
	userId, errCode := r.jwtService.Validate(request.Token)
	if errCode != nil {
		return "", errCode
	}

	// get user
	var user models.User
	err := r.userRepository.Get(&user, userId)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return "", wrapper.UnauthorizedError(errors.New("not authorized"))
	}

	return r.jwtService.CreateToken(user)
}
