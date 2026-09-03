package song

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

type SaveImageToSong struct {
	repository              repository.SongRepository
	storageFilePathProvider provider.StorageFilePathProvider
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewSaveImageToSong(
	repository repository.SongRepository,
	storageFilePathProvider provider.StorageFilePathProvider,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) SaveImageToSong {
	return SaveImageToSong{
		repository:              repository,
		storageFilePathProvider: storageFilePathProvider,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (s SaveImageToSong) Handle(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	var song model.Song
	if err := s.repository.Get(&song, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	if song.ImageURL != nil {
		if errCode := s.storageService.DeleteFile(*song.ImageURL); errCode != nil {
			return errCode
		}
	}

	song.UpdatedAt = time.Now().UTC()
	imagePath := s.storageFilePathProvider.GetSongImagePath(file, song)

	if errCode := s.storageService.Upload(file, imagePath); errCode != nil {
		return errCode
	}

	song.ImageURL = (*internal.FilePath)(&imagePath)
	if err := s.repository.Update(&song); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := s.messagePublisherService.Publish(topics.SongsUpdatedTopic, []uuid.UUID{song.ID}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
