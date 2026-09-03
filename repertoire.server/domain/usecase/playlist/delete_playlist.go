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

type DeletePlaylist struct {
	playlistRepository      repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewDeletePlaylist(
	playlistRepository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) DeletePlaylist {
	return DeletePlaylist{
		playlistRepository:      playlistRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeletePlaylist) Handle(id uuid.UUID) *httperror.ErrorCode {
	var playlist model.Playlist
	if err := d.playlistRepository.Get(&playlist, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return httperror.NotFoundError(errors.New("playlist not found"))
	}

	if err := d.playlistRepository.Delete([]uuid.UUID{id}); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.PlaylistsDeletedTopic, []model.Playlist{playlist}); err != nil {
		return httperror.MessagePublisherError(err)
	}
	return nil
}
