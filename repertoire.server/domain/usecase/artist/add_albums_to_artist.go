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

type AddAlbumsToArtist struct {
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewAddAlbumsToArtist(
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) AddAlbumsToArtist {
	return AddAlbumsToArtist{
		albumRepository:         albumRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (a AddAlbumsToArtist) Handle(request requests.AddAlbumsToArtistRequest) *httperror.ErrorCode {
	var albums []model.Album
	if err := a.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs); err != nil {
		return httperror.DatabaseError(err)
	}

	for i, album := range albums {
		if album.ArtistID != nil {
			return httperror.ConflictError(errors.New("album " + album.ID.String() + " already has an artist"))
		}

		// update the whole album's artist
		albums[i].ArtistID = &request.ID
		for j := range album.Songs {
			albums[i].Songs[j].ArtistID = &request.ID
		}
	}

	if err := a.albumRepository.UpdateAllWithSongs(&albums); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := a.messagePublisherService.Publish(topics.AlbumsUpdatedTopic, request.AlbumIDs); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}
