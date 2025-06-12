package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
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

func (r RemoveAlbumsFromArtist) Handle(request requests.RemoveAlbumsFromArtistRequest) *wrapper.ErrorCode {
	var albums []model.Album
	err := r.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for i, album := range albums {
		if album.ArtistID == nil || *album.ArtistID != request.ID {
			return wrapper.ConflictError(errors.New("album " + album.ID.String() + " is not owned by this artist"))
		}

		albums[i].ArtistID = nil
		for j := range album.Songs {
			albums[i].Songs[j].ArtistID = nil
		}
	}

	err = r.albumRepository.UpdateAllWithSongs(&albums)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = r.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, request.AlbumIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
