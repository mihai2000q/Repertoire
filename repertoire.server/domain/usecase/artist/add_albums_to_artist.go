package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddAlbumsToArtist struct {
	albumRepository repository.AlbumRepository
}

func NewAddAlbumsToArtist(albumRepository repository.AlbumRepository) AddAlbumsToArtist {
	return AddAlbumsToArtist{albumRepository: albumRepository}
}

func (a AddAlbumsToArtist) Handle(request requests.AddAlbumsToArtistRequest) *wrapper.ErrorCode {
	var albums []model.Album
	err := a.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for i, album := range albums {
		if album.ArtistID != nil {
			return wrapper.BadRequestError(errors.New("album " + album.ID.String() + " already has an artist"))
		}

		// update the whole album's artist
		albums[i].ArtistID = &request.ID
		for j := range album.Songs {
			albums[i].Songs[j].ArtistID = &request.ID
		}
	}

	err = a.albumRepository.UpdateAllWithSongs(&albums)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
