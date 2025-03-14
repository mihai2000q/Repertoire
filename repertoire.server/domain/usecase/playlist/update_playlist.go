package playlist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdatePlaylist struct {
	repository              repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewUpdatePlaylist(
	repository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) UpdatePlaylist {
	return UpdatePlaylist{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (u UpdatePlaylist) Handle(request requests.UpdatePlaylistRequest) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := u.repository.Get(&playlist, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return wrapper.NotFoundError(errors.New("playlist not found"))
	}

	playlist.Title = request.Title
	playlist.Description = request.Description

	err = u.repository.Update(&playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = u.messagePublisherService.Publish(topics.PlaylistUpdatedTopic, playlist)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
