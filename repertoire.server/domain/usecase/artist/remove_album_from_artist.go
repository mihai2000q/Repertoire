package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type RemoveAlbumsFromArtist struct {
	albumRepository repository.AlbumRepository
}

func NewRemoveAlbumsFromArtist(albumRepository repository.AlbumRepository) RemoveAlbumsFromArtist {
	return RemoveAlbumsFromArtist{albumRepository: albumRepository}
}

func (r RemoveAlbumsFromArtist) Handle(request requests.RemoveAlbumsFromArtistRequest) *wrapper.ErrorCode {
	var albums []model.Album
	err := r.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for _, album := range albums {
		if album.ArtistID == nil || *album.ArtistID != request.ID {
			return wrapper.BadRequestError(errors.New("album " + album.ID.String() + " is not owned by this artist"))
		}

		album.ArtistID = nil
		for i := range album.Songs {
			album.Songs[i].ArtistID = nil
		}

		err = r.albumRepository.UpdateWithAssociations(&album)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	return nil
}
