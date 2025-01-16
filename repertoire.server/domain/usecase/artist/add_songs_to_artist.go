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

	for i, song := range songs {
		if song.ArtistID != nil {
			return wrapper.BadRequestError(errors.New("song " + song.ID.String() + "already has an artist"))
		}

		// update the whole album's artist, including the other songs
		if song.Album != nil {
			songs[i].Album.ArtistID = &request.ID
			for j := range song.Album.Songs {
				songs[i].Album.Songs[j].ArtistID = &request.ID
			}
		} else {
			songs[i].ArtistID = &request.ID
		}
	}

	err = a.songRepository.UpdateAllWithAssociations(&songs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
