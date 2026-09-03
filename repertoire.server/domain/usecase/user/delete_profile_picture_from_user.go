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
	userRepository repository.UserRepository
	jwtService     service.JwtService
	storageService service.StorageService
}

func NewDeleteProfilePictureFromUser(
	userRepository repository.UserRepository,
	jwtService service.JwtService,
	storageService service.StorageService,
) DeleteProfilePictureFromUser {
	return DeleteProfilePictureFromUser{
		userRepository: userRepository,
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
	if err := d.userRepository.Get(&user, id); err != nil {
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
	if err := d.userRepository.Update(&user); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
