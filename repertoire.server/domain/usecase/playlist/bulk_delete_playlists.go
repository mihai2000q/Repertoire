package playlist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type BulkDeletePlaylists struct {
	repository              repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewBulkDeletePlaylists(
	repository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) BulkDeletePlaylists {
	return BulkDeletePlaylists{
		repository:              repository,
		messagePublisherService: messagePublisherService,
	}
}

func (b BulkDeletePlaylists) Handle(request requests.BulkDeletePlaylistsRequest) *httperror.ErrorCode {
	var playlists []model.Playlist
	if err := b.repository.GetAllByIDs(&playlists, request.IDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(playlists) == 0 {
		return httperror.NotFoundError(errors.New("playlists not found"))
	}

	if err := b.repository.Delete(request.IDs); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := b.messagePublisherService.Publish(topics.PlaylistsDeletedTopic, playlists); err != nil {
		return httperror.MessagePublisherError(err)
	}
	return nil
}
