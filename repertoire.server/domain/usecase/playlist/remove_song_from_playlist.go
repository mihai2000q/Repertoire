package playlist

import (
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type RemoveSongFromPlaylist struct {
	songRepository repository.SongRepository
	repository     repository.PlaylistRepository
}

func NewRemoveSongFromPlaylist(
	songRepository repository.SongRepository,
	repository repository.PlaylistRepository,
) RemoveSongFromPlaylist {
	return RemoveSongFromPlaylist{
		songRepository: songRepository,
		repository:     repository,
	}
}

func (r RemoveSongFromPlaylist) Handle(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var playlist model.Playlist
	err := r.repository.Get(&playlist, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var song model.Song
	err = r.songRepository.Get(&song, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = r.repository.RemoveSong(&playlist, &song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
