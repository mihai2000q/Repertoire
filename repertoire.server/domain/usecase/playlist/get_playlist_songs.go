package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetPlaylistSongs struct {
	repository repository.PlaylistRepository
}

func NewGetPlaylistSongs(repository repository.PlaylistRepository) GetPlaylistSongs {
	return GetPlaylistSongs{
		repository: repository,
	}
}

func (g GetPlaylistSongs) Handle(request requests.GetPlaylistSongsRequest) (songs []model.Song, e *wrapper.ErrorCode) {
	if len(request.OrderBy) == 0 {
		request.OrderBy = []string{"song_track_no"}
	}

	var playlistSongs []model.PlaylistSong
	err := g.repository.GetPlaylistSongsWithSongs(
		&playlistSongs,
		request.ID,
		request.CurrentPage,
		request.PageSize,
		request.OrderBy,
	)
	if err != nil {
		return songs, wrapper.InternalServerError(err)
	}

	for _, playlistSong := range playlistSongs {
		songs = append(songs, g.mapToSong(playlistSong))
	}

	return songs, nil
}

func (g GetPlaylistSongs) mapToSong(playlistSong model.PlaylistSong) model.Song {
	song := playlistSong.Song

	song.PlaylistSongID = playlistSong.ID
	song.PlaylistTrackNo = playlistSong.SongTrackNo
	song.PlaylistCreatedAt = playlistSong.CreatedAt

	song.ToFullImageURL()

	return song
}
