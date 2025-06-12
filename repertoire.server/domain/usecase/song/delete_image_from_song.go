package song

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteImageFromSong struct {
	repository              repository.SongRepository
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewDeleteImageFromSong(
	repository repository.SongRepository,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,
) DeleteImageFromSong {
	return DeleteImageFromSong{
		repository:              repository,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteImageFromSong) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.repository.Get(&song, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}
	if song.ImageURL == nil {
		return wrapper.ConflictError(errors.New("song does not have an image"))
	}

	errCode := d.storageService.DeleteFile(*song.ImageURL)
	if errCode != nil {
		return errCode
	}

	song.ImageURL = nil
	err = d.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.SongsUpdatedTopic, []uuid.UUID{song.ID})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
