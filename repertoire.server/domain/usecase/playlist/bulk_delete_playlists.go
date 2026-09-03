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
	playlistRepository      repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewBulkDeletePlaylists(
	playlistRepository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) BulkDeletePlaylists {
	return BulkDeletePlaylists{
		playlistRepository:      playlistRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (b BulkDeletePlaylists) Handle(request requests.BulkDeletePlaylistsRequest) *httperror.ErrorCode {
	var playlists []model.Playlist
	if err := b.playlistRepository.GetAllByIDs(&playlists, request.IDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(playlists) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("playlists not found"))
	}

	if err := b.playlistRepository.Delete(request.IDs); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := b.messagePublisherService.Publish(topics.PlaylistsDeletedTopic, playlists); err != nil {
		return httperror.MessagePublisherError(err)
	}
	return nil
}
