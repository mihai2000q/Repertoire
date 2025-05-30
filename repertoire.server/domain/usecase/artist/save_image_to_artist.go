package artist

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
	"time"

	"github.com/google/uuid"
)

type SaveImageToArtist struct {
	repository              repository.ArtistRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewSaveImageToArtist(
	repository repository.ArtistRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) SaveImageToArtist {
	return SaveImageToArtist{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (s SaveImageToArtist) Handle(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	var artist model.Artist
	err := s.repository.Get(&artist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return wrapper.NotFoundError(errors.New("artist not found"))
	}

	if artist.ImageURL != nil {
		errCode := s.storageService.DeleteFile(*artist.ImageURL)
		if errCode != nil {
			return errCode
		}
	}

	artist.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetArtistImagePath(file, artist)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	artist.ImageURL = (*internal.FilePath)(&imagePath)
	err = s.repository.Update(&artist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = s.messagePublisherService.Publish(topics.ArtistUpdatedTopic, artist.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
