package provider

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
	"repertoire/utils"
)

type CurrentUserProvider interface {
	Get(token string) (user models.User, e *utils.ErrorCode)
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

func (c *currentUserProvider) Get(token string) (user models.User, e *utils.ErrorCode) {
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
