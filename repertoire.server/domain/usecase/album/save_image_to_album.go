package album

import (
	"errors"
	"mime/multipart"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"time"

	"github.com/google/uuid"
)

type SaveImageToAlbum struct {
	albumRepository         repository.AlbumRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewSaveImageToAlbum(
	albumRepository repository.AlbumRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) SaveImageToAlbum {
	return SaveImageToAlbum{
		albumRepository:         albumRepository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (s SaveImageToAlbum) Handle(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	var album model.Album
	if err := s.albumRepository.Get(&album, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return httperror.NotFoundError(errors.New("album not found"))
	}

	if album.ImageURL != nil {
		if errCode := s.storageService.DeleteFile(*album.ImageURL); errCode != nil {
			return errCode
		}
	}

	album.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetAlbumImagePath(file, album)

	if errCode := s.storageService.Upload(file, imagePath); errCode != nil {
		return errCode
	}

	album.ImageURL = (*internal.FilePath)(&imagePath)
	if err := s.albumRepository.Update(&album); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := s.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, []uuid.UUID{album.ID}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
