package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongsToArtist struct {
	songRepository repository.SongRepository
}

func NewAddSongsToArtist(songRepository repository.SongRepository) AddSongsToArtist {
	return AddSongsToArtist{songRepository: songRepository}
}

func (a AddSongsToArtist) Handle(request requests.AddSongsToArtistRequest) *wrapper.ErrorCode {
	var songs []model.Song
	err := a.songRepository.GetAllByIDsWithSongs(&songs, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for _, song := range songs {
		if song.ArtistID != nil {
			return wrapper.BadRequestError(errors.New("song " + song.ID.String() + "already has an artist"))
		}

		// update the whole album's artist, including the other songs
		if song.Album != nil {
			song.Album.ArtistID = &request.ID
			for i := range song.Album.Songs {
				song.Album.Songs[i].ArtistID = &request.ID
			}
		} else {
			song.ArtistID = &request.ID
		}

		err = a.songRepository.UpdateWithAssociations(&song)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	return nil
}
