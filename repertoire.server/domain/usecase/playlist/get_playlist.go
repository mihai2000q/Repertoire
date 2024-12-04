package playlist

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type GetPlaylist struct {
	repository repository.PlaylistRepository
}

func NewGetPlaylist(repository repository.PlaylistRepository) GetPlaylist {
	return GetPlaylist{
		repository: repository,
	}
}

func (g GetPlaylist) Handle(id uuid.UUID) (playlist model.Playlist, e *wrapper.ErrorCode) {
	err := g.repository.GetWithAssociations(&playlist, id)
	if err != nil {
		return playlist, wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(playlist).IsZero() {
		return playlist, wrapper.NotFoundError(errors.New("playlist not found"))
	}
	g.LoadSongs(&playlist)
	return playlist, nil
}

func (g GetPlaylist) LoadSongs(playlist *model.Playlist) {
	for _, playlistSong := range playlist.PlaylistSongs {
		newSong := playlistSong.Song
		newSong.PlaylistTrackNo = playlistSong.SongTrackNo
		newSong.PlaylistCreatedAt = playlistSong.CreatedAt
		playlist.Songs = append(playlist.Songs, newSong)
	}
}
