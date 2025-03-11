package album

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type SaveImageToAlbum struct {
	repository              repository.AlbumRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewSaveImageToAlbum(
	repository repository.AlbumRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) SaveImageToAlbum {
	return SaveImageToAlbum{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (s SaveImageToAlbum) Handle(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	var album model.Album
	err := s.repository.Get(&album, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	imagePath := s.storageFilePathProvider.GetAlbumImagePath(file, album)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	album.ImageURL = (*internal.FilePath)(&imagePath)
	err = s.repository.Update(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = s.messagePublisherService.Publish(topics.AlbumUpdatedTopic, album.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
