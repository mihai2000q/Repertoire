package playlist

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveSongFromPlaylist struct {
	repository repository.PlaylistRepository
}

func NewMoveSongFromPlaylist(repository repository.PlaylistRepository) MoveSongFromPlaylist {
	return MoveSongFromPlaylist{repository: repository}
}

func (m MoveSongFromPlaylist) Handle(request requests.MoveSongFromPlaylistRequest) *wrapper.ErrorCode {
	var playlistSongs []model.PlaylistSong
	err := m.repository.GetPlaylistSongs(&playlistSongs, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	index, overIndex, err := m.getIndexes(playlistSongs, request.SongID, request.OverSongID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	playlistSongs = m.move(playlistSongs, index, overIndex)

	err = m.repository.UpdateAllPlaylistSongs(&playlistSongs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (MoveSongFromPlaylist) getIndexes(
	playlistSongs []model.PlaylistSong,
	id uuid.UUID,
	overID uuid.UUID,
) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(playlistSongs); i++ {
		if playlistSongs[i].SongID == id {
			index = &i
		} else if playlistSongs[i].SongID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("song not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over song not found")
	}

	return *index, *overIndex, nil
}

func (MoveSongFromPlaylist) move(playlistSongs []model.PlaylistSong, index int, overIndex int) []model.PlaylistSong {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			playlistSongs[i].SongTrackNo = playlistSongs[i].SongTrackNo - 1
		}
	} else {
		for i := overIndex; i <= index; i++ {
			playlistSongs[i].SongTrackNo = playlistSongs[i].SongTrackNo + 1
		}
	}
	playlistSongs[index].SongTrackNo = uint(overIndex) + 1

	return playlistSongs
}
