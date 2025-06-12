package playlist

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

type DeleteImageFromPlaylist struct {
	repository              repository.PlaylistRepository
	storageService          service.StorageService
	messagePublisherService service.MessagePublisherService
}

func NewDeleteImageFromPlaylist(
	repository repository.PlaylistRepository,
	storageService service.StorageService,
	messagePublisherService service.MessagePublisherService,

) DeleteImageFromPlaylist {
	return DeleteImageFromPlaylist{
		repository:              repository,
		storageService:          storageService,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteImageFromPlaylist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := d.repository.Get(&playlist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}
	if playlist.ImageURL == nil {
		return wrapper.ConflictError(errors.New("playlist does not have an image"))
	}

	errCode := d.storageService.DeleteFile(*playlist.ImageURL)
	if errCode != nil {
		return errCode
	}

	playlist.ImageURL = nil
	err = d.repository.Update(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.PlaylistUpdatedTopic, playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
