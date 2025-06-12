package user

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type DeleteProfilePictureFromUser struct {
	repository     repository.UserRepository
	jwtService     service.JwtService
	storageService service.StorageService
}

func NewDeleteProfilePictureFromUser(
	repository repository.UserRepository,
	jwtService service.JwtService,
	storageService service.StorageService,
) DeleteProfilePictureFromUser {
	return DeleteProfilePictureFromUser{
		repository:     repository,
		jwtService:     jwtService,
		storageService: storageService,
	}
}

func (d DeleteProfilePictureFromUser) Handle(token string) *wrapper.ErrorCode {
	id, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var user model.User
	err := d.repository.Get(&user, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return wrapper.NotFoundError(errors.New("user not found"))
	}
	if user.ProfilePictureURL == nil {
		return wrapper.ConflictError(errors.New("user does not have a profile picture"))
	}

	errCode = d.storageService.DeleteFile(*user.ProfilePictureURL)
	if errCode != nil {
		return errCode
	}

	user.ProfilePictureURL = nil
	err = d.repository.Update(&user)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
