package song

import (
	"math/rand"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type ShufflePlaylistSongs struct {
	repository repository.PlaylistRepository
}

func NewShufflePlaylistSongs(repository repository.PlaylistRepository) ShufflePlaylistSongs {
	return ShufflePlaylistSongs{repository: repository}
}

func (s ShufflePlaylistSongs) Handle(request requests.ShufflePlaylistSongsRequest) *httperror.ErrorCode {
	var songs []model.PlaylistSong
	if err := s.repository.GetPlaylistSongs(&songs, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}

	for i := range songs {
		j := rand.Intn(i + 1)
		songs[i], songs[j] = songs[j], songs[i]
		songs[i].SongTrackNo = uint(i + 1)
		songs[j].SongTrackNo = uint(j + 1)
	}

	if err := s.repository.UpdateAllPlaylistSongs(&songs); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
