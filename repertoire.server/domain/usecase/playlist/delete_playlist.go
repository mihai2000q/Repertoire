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

type DeletePlaylist struct {
	repository              repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewDeletePlaylist(
	repository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) DeletePlaylist {
	return DeletePlaylist{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeletePlaylist) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := d.repository.Get(&playlist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}

	err = d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.PlaylistsDeletedTopic, []model.Playlist{playlist})
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
