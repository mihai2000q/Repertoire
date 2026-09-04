package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
)

type RemoveAlbumsFromArtist struct {
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewRemoveAlbumsFromArtist(
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) RemoveAlbumsFromArtist {
	return RemoveAlbumsFromArtist{
		albumRepository:         albumRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (r RemoveAlbumsFromArtist) Handle(request requests.RemoveAlbumsFromArtistRequest) *httperror.ErrorCode {
	var albums []model.Album
	if err := r.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(albums) != len(request.AlbumIDs) {
		return httperror.NotFoundError(errors.New("albums not found"))
	}

	for i, album := range albums {
		if album.ArtistID == nil || *album.ArtistID != request.ID {
			return httperror.ConflictError(errors.New("album " + album.ID.String() + " is not owned by this artist"))
		}

		albums[i].ArtistID = nil
		for j := range album.Songs {
			albums[i].Songs[j].ArtistID = nil
		}
	}

	if err := r.albumRepository.UpdateAllWithSongs(&albums); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := r.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, request.AlbumIDs); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
