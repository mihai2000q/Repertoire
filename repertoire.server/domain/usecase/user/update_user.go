package user

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateUser struct {
	repository repository.UserRepository
	jwtService service.JwtService
}

func NewUpdateUser(
	repository repository.UserRepository,
	jwtService service.JwtService,
) UpdateUser {
	return UpdateUser{
		repository: repository,
		jwtService: jwtService,
	}
}

func (u UpdateUser) Handle(request requests.UpdateUserRequest, token string) *wrapper.ErrorCode {
	id, errCode := u.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var user model.User
	err := u.repository.Get(&user, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return wrapper.NotFoundError(errors.New("user not found"))
	}

	user.Name = request.Name

	err = u.repository.Update(&user)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
