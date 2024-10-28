package provider

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/model"
	"repertoire/utils/wrapper"
)

type CurrentUserProvider interface {
	Get(token string) (user model.User, e *wrapper.ErrorCode)
}

type currentUserProvider struct {
	jwtService     service.JwtService
	userRepository repository.UserRepository
}

func NewCurrentUserProvider(
	jwtService service.JwtService,
	userRepository repository.UserRepository,
) CurrentUserProvider {
	return &currentUserProvider{
		jwtService:     jwtService,
		userRepository: userRepository,
	}
}

func (c *currentUserProvider) Get(token string) (user model.User, e *wrapper.ErrorCode) {
	userID, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return user, errCode
	}

	err := c.userRepository.Get(&user, userID)
	if err != nil {
		return user, wrapper.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return user, wrapper.NotFoundError(errors.New("user not found"))
	}

	return user, nil
}
