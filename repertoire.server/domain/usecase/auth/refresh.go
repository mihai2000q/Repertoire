package auth

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
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
	userID, errCode := r.jwtService.Validate(request.Token)
	if errCode != nil {
		return "", errCode
	}

	// get user
	var user model.User
	err := r.userRepository.Get(&user, userID)
	if err != nil {
		return "", wrapper.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return "", wrapper.UnauthorizedError(errors.New("not authorized"))
	}

	return r.jwtService.CreateToken(user)
}
