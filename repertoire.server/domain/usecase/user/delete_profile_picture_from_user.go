package user

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
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

func (d DeleteProfilePictureFromUser) Handle(token string) *httperror.ErrorCode {
	id, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var user model.User
	if err := d.repository.Get(&user, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return httperror.NotFoundError(errors.New("user not found"))
	}
	if user.ProfilePictureURL == nil {
		return nil
	}

	if errCode = d.storageService.DeleteFile(*user.ProfilePictureURL); errCode != nil {
		return errCode
	}

	user.ProfilePictureURL = nil
	if err := d.repository.Update(&user); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
