package playlist

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

func (d DeleteImageFromPlaylist) Handle(id uuid.UUID) *httperror.ErrorCode {
	var playlist model.Playlist
	if err := d.repository.Get(&playlist, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return httperror.NotFoundError(errors.New("playlist not found"))
	}
	if playlist.ImageURL == nil {
		return httperror.ConflictError(errors.New("playlist does not have an image"))
	}

	errCode := d.storageService.DeleteFile(*playlist.ImageURL)
	if errCode != nil {
		return errCode
	}

	playlist.ImageURL = nil
	if err := d.repository.Update(&playlist); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.PlaylistUpdatedTopic, playlist); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
