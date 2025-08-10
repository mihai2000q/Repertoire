package playlist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
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

func (b BulkDeletePlaylists) Handle(request requests.BulkDeletePlaylistsRequest) *wrapper.ErrorCode {
	var playlists []model.Playlist
	err := b.repository.GetAllByIDs(&playlists, request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(playlists) == 0 {
		return wrapper.NotFoundError(errors.New("playlists not found"))
	}

	err = b.repository.Delete(request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = b.messagePublisherService.Publish(topics.PlaylistsDeletedTopic, playlists)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
