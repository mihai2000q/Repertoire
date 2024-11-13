package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type RemoveSongFromArtist struct {
	songRepository repository.SongRepository
}

func NewRemoveSongFromArtist(songRepository repository.SongRepository) RemoveSongFromArtist {
	return RemoveSongFromArtist{songRepository: songRepository}
}

func (r RemoveSongFromArtist) Handle(request requests.RemoveSongsFromArtistRequest) *wrapper.ErrorCode {
	var songs []model.Song
	err := r.songRepository.GetAllByIDs(&songs, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for i, song := range songs {
		if song.ArtistID == nil || *song.ArtistID != request.ID {
			return wrapper.BadRequestError(errors.New("song " + song.ID.String() + " is not owned by this artist"))
		}

		songs[i].ArtistID = nil
	}

	err = r.songRepository.UpdateAll(&songs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
