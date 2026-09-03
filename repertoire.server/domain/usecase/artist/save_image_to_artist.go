package artist

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

func (s SaveImageToArtist) Handle(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	var artist model.Artist
	if err := s.repository.Get(&artist, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(artist).IsZero() {
		return httperror.NotFoundError(errors.New("artist not found"))
	}

	if artist.ImageURL != nil {
		if errCode := s.storageService.DeleteFile(*artist.ImageURL); errCode != nil {
			return errCode
		}
	}

	artist.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetArtistImagePath(file, artist)

	if errCode := s.storageService.Upload(file, imagePath); errCode != nil {
		return errCode
	}

	artist.ImageURL = (*internal.FilePath)(&imagePath)
	if err := s.repository.Update(&artist); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := s.messagePublisherService.Publish(topics.ArtistUpdatedTopic, artist.ID); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
