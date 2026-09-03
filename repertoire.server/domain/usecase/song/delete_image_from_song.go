package song

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
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

func (d DeleteImageFromSong) Handle(id uuid.UUID) *httperror.ErrorCode {
	var song model.Song
	if err := d.repository.Get(&song, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}
	if song.ImageURL == nil {
		return nil
	}

	if errCode := d.storageService.DeleteFile(*song.ImageURL); errCode != nil {
		return errCode
	}

	song.ImageURL = nil
	if err := d.repository.Update(&song); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.SongsUpdatedTopic, []uuid.UUID{song.ID}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
