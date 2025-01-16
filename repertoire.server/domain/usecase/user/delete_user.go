package user

import (
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/wrapper"
)

type DeleteUser struct {
	repository              repository.UserRepository
	jwtService              service.JwtService
	storageService          service.StorageService
	storageFilePathProvider provider.StorageFilePathProvider
}

func NewDeleteUser(
	repository repository.UserRepository,
	jwtService service.JwtService,
	storageService service.StorageService,
	storageFilePathProvider provider.StorageFilePathProvider,
) DeleteUser {
	return DeleteUser{
		repository:              repository,
		jwtService:              jwtService,
		storageService:          storageService,
		storageFilePathProvider: storageFilePathProvider,
	}
}

func (d DeleteUser) Handle(token string) *wrapper.ErrorCode {
	id, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	directoryPath := d.storageFilePathProvider.GetUserDirectoryPath(id)
	err := d.storageService.DeleteDirectory(directoryPath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
