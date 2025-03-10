package song

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

func (s SaveImageToSong) Handle(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := s.repository.Get(&song, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	imagePath := s.storageFilePathProvider.GetSongImagePath(file, song)

	err = s.storageService.Upload(file, imagePath)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	song.ImageURL = (*internal.FilePath)(&imagePath)
	err = s.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = s.messagePublisherService.Publish(topics.SongsUpdatedTopic, []uuid.UUID{song.ID})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
