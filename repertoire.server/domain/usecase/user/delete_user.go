package user

import (
	"net/http"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
)

type DeleteUser struct {
	repository              repository.UserRepository
	jwtService              service.JwtService
	storageService          service.StorageService
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewDeleteUser(
	repository repository.UserRepository,
	jwtService service.JwtService,
	storageService service.StorageService,
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) DeleteUser {
	return DeleteUser{
		repository:              repository,
		jwtService:              jwtService,
		storageService:          storageService,
		storageFilePathProvider: storageFilePathProvider,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteUser) Handle(token string) *wrapper.ErrorCode {
	id, errCode := d.jwtService.GetUserIdFromJwt(token)
	if errCode != nil {
		return errCode
	}

	directoryPath := d.storageFilePathProvider.GetUserDirectoryPath(id)
	errCode = d.storageService.DeleteDirectory(directoryPath)
	if errCode != nil && errCode.Code != http.StatusNotFound {
		return errCode
	}

	err := d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.UserDeletedTopic, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
