package album

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongToAlbum struct {
	albumRepository repository.AlbumRepository
	repository      repository.SongRepository
}

func NewAddSongToAlbum(
	albumRepository repository.AlbumRepository,
	repository repository.SongRepository,
) AddSongToAlbum {
	return AddSongToAlbum{
		albumRepository: albumRepository,
		repository:      repository,
	}
}

func (a AddSongToAlbum) Handle(request requests.AddSongToAlbumRequest) *wrapper.ErrorCode {
	var song model.Song
	err := a.repository.Get(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var count int64
	err = a.albumRepository.CountSongs(&count, &request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	song.AlbumID = &request.ID
	trackNo := uint(count) + 1
	song.AlbumTrackNo = &trackNo

	err = a.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
