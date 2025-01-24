package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
)

type AddAlbumsToPlaylist struct {
	repository      repository.PlaylistRepository
	albumRepository repository.AlbumRepository
}

func NewAddAlbumsToPlaylist(
	repository repository.PlaylistRepository,
	albumRepository repository.AlbumRepository,
) AddAlbumsToPlaylist {
	return AddAlbumsToPlaylist{
		repository:      repository,
		albumRepository: albumRepository,
	}
}

func (a AddAlbumsToPlaylist) Handle(request requests.AddAlbumsToPlaylistRequest) *wrapper.ErrorCode {
	var oldPlaylistSongs []model.PlaylistSong
	err := a.repository.GetPlaylistSongs(&oldPlaylistSongs, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var albums []model.Album
	err = a.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	songsLength := len(oldPlaylistSongs) + 1
	currentTrackNo := uint(songsLength)
	var playlistSongs []model.PlaylistSong
	for _, album := range albums {
		for _, song := range album.Songs {
			if slices.ContainsFunc(oldPlaylistSongs, func(p model.PlaylistSong) bool {
				return p.SongID == song.ID
			}) {
				continue
			}

			playlistSong := model.PlaylistSong{
				PlaylistID:  request.ID,
				SongID:      song.ID,
				SongTrackNo: currentTrackNo,
			}
			playlistSongs = append(playlistSongs, playlistSong)
			currentTrackNo++
		}
	}

	err = a.repository.AddSongs(&playlistSongs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
