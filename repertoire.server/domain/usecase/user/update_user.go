package user

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type UpdateUser struct {
	userRepository repository.UserRepository
	jwtService     service.JwtService
}

func NewUpdateUser(
	userRepository repository.UserRepository,
	jwtService service.JwtService,
) UpdateUser {
	return UpdateUser{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (u UpdateUser) Handle(request requests.UpdateUserRequest, token string) *httperror.ErrorCode {
	id, errCode := u.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var user model.User
	if err := u.userRepository.Get(&user, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return httperror.NotFoundError(errors.New("user not found"))
	}

	user.Name = request.Name

	if err := u.userRepository.Update(&user); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
