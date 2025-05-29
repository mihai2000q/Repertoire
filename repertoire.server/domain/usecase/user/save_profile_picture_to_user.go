package user

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/wrapper"
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

func (s SaveProfilePictureToUser) Handle(file *multipart.FileHeader, token string) *wrapper.ErrorCode {
	id, errCode := s.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	var user model.User
	err := s.repository.Get(&user, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(user).IsZero() {
		return wrapper.NotFoundError(errors.New("user not found"))
	}

	if user.ProfilePictureURL != nil {
		errCode = s.storageService.DeleteFile(*user.ProfilePictureURL)
		if errCode != nil {
			return errCode
		}
	}

	user.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetUserProfilePicturePath(file, user)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	user.ProfilePictureURL = (*internal.FilePath)(&imagePath)
	err = s.repository.Update(&user)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
