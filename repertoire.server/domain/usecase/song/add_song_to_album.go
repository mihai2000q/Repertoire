package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongToAlbum struct {
	repository repository.SongRepository
}

func NewAddSongToAlbum(repository repository.SongRepository) AddSongToAlbum {
	return AddSongToAlbum{repository: repository}
}

func (a AddSongToAlbum) Handle(request requests.AddSongToAlbumRequest) *wrapper.ErrorCode {
	var song model.Song
	err := a.repository.Get(&song, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var count int64
	err = a.repository.CountByAlbum(&count, &request.AlbumID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	song.AlbumID = &request.AlbumID
	trackNo := uint(count) + 1
	song.AlbumTrackNo = &trackNo

	err = a.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
