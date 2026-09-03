package user

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"time"
)

type SaveProfilePictureToUser struct {
	repository              repository.UserRepository
	storageFilePathProvider provider.StorageFilePathProvider
	jwtService              service.JwtService
	storageService          service.StorageService
}

func NewSaveProfilePictureToUser(
	repository repository.UserRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	jwtService service.JwtService,
	storageService service.StorageService,
) SaveProfilePictureToUser {
	return SaveProfilePictureToUser{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		jwtService:              jwtService,
		storageService:          storageService,
	}
}

func (s SaveProfilePictureToUser) Handle(file *multipart.FileHeader, token string) *httperror.ErrorCode {
	id, errCode := s.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var user model.User
	if err := s.repository.Get(&user, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return httperror.NotFoundError(errors.New("user not found"))
	}

	if user.ProfilePictureURL != nil {
		if errCode = s.storageService.DeleteFile(*user.ProfilePictureURL); errCode != nil {
			return errCode
		}
	}

	user.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetUserProfilePicturePath(file, user)

	if errCode := s.storageService.Upload(file, imagePath); errCode != nil {
		return errCode
	}

	user.ProfilePictureURL = (*internal.FilePath)(&imagePath)
	if err := s.repository.Update(&user); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
