package album

import (
	"errors"
	"net/http"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type DeleteAlbum struct {
	repository              repository.AlbumRepository
	storageService          service.StorageService
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewDeleteAlbum(
	repository repository.AlbumRepository,
	storageService service.StorageService,
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) DeleteAlbum {
	return DeleteAlbum{
		repository:              repository,
		storageService:          storageService,
		storageFilePathProvider: storageFilePathProvider,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteAlbum) Handle(request requests.DeleteAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	var err error
	if request.WithSongs {
		err = d.repository.GetWithSongs(&album, request.ID)
	} else {
		err = d.repository.Get(&album, request.ID)
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	directoryPath := d.storageFilePathProvider.GetAlbumDirectoryPath(album)
	errCode := d.storageService.DeleteDirectory(directoryPath)
	if errCode != nil && errCode.Code != http.StatusNotFound {
		return errCode
	}

	if request.WithSongs {
		err = d.repository.DeleteWithSongs(request.ID)
	} else {
		err = d.repository.Delete(request.ID)
	}
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.AlbumDeletedTopic, album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
