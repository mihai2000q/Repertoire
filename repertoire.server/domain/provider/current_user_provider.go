package provider

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/config"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
)

type CurrentUserProvider struct {
	jwtService     service.JwtService
	userRepository repository.UserRepository
	env            config.Env
}

func NewCurrentUserProvider(
	jwtService service.JwtService,
	userRepository repository.UserRepository,
	env config.Env,
) CurrentUserProvider {
	return CurrentUserProvider{
		jwtService:     jwtService,
		userRepository: userRepository,
		env:            env,
	}
}

func (c *CurrentUserProvider) Get(token string) (user models.User, e *utils.ErrorCode) {
	userId, errCode := c.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return user, errCode
	}

	err := c.userRepository.Get(&user, userId)
	if err != nil {
		return user, utils.InternalServerError(err)
	}
	if user.ID == uuid.Nil {
		return user, utils.NotFoundError(errors.New("user not found"))
	}

	return user, nil
}
