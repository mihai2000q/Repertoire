package song

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
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

func (g GetPlaylistSongs) Handle(request requests.GetPlaylistSongsRequest) (result pagination.WithTotalCount[model.Song], e *httperror.ErrorCode) {
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
		return result, httperror.DatabaseError(err)
	}

	err = g.repository.GetPlaylistSongsCount(&result.TotalCount, request.ID)
	if err != nil {
		return result, httperror.DatabaseError(err)
	}

	var songs []model.Song
	for _, playlistSong := range playlistSongs {
		songs = append(songs, g.mapToSong(playlistSong))
	}
	result.Models = songs

	return result, nil
}

func (g GetPlaylistSongs) mapToSong(playlistSong model.PlaylistSong) model.Song {
	song := playlistSong.Song

	song.PlaylistSongID = playlistSong.ID
	song.PlaylistTrackNo = playlistSong.SongTrackNo
	song.PlaylistCreatedAt = playlistSong.CreatedAt

	song.ToFullImageURL()

	return song
}
