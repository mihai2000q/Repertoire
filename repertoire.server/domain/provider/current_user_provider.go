package provider

import (
	"repertoire/config"
	"repertoire/data/repository"
	"repertoire/data/service"
	"repertoire/models"
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

func (c *CurrentUserProvider) Get(token string) (*models.User, error) {
	userId, err := c.jwtService.GetUserIdFromJwt(token)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = c.userRepository.Get(&user, userId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
