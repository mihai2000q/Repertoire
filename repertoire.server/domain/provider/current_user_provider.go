package provider

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type CurrentUserProvider interface {
	Get(token string) (model.User, *wrapper.ErrorCode)
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
	if reflect.ValueOf(user).IsZero() {
		return user, wrapper.NotFoundError(errors.New("user not found"))
	}

	return user, nil
}
