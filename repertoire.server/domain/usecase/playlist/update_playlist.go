package playlist

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
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

func (u UpdatePlaylist) Handle(request requests.UpdatePlaylistRequest) *httperror.ErrorCode {
	var playlist model.Playlist
	if err := u.repository.Get(&playlist, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return httperror.NotFoundError(errors.New("playlist not found"))
	}

	playlist.Title = request.Title
	playlist.Description = request.Description

	if err := u.repository.Update(&playlist); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := u.messagePublisherService.Publish(topics.PlaylistUpdatedTopic, playlist); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
