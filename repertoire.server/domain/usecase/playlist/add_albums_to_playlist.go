package playlist

import (
	"repertoire/server/api/requests"
	"repertoire/server/api/responses"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
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

func (a AddAlbumsToPlaylist) Handle(
	request requests.AddAlbumsToPlaylistRequest,
) (*responses.AddAlbumsToPlaylistResponse, *wrapper.ErrorCode) {
	var playlistSongs []model.PlaylistSong
	err := a.repository.GetPlaylistSongs(&playlistSongs, request.ID)
	if err != nil {
		return nil, wrapper.InternalServerError(err)
	}

	var albums []model.Album
	err = a.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs)
	if err != nil {
		return nil, wrapper.InternalServerError(err)
	}

	var duplicateSongIDs []uuid.UUID
	var duplicateAlbumIDs []uuid.UUID

	songsLength := len(playlistSongs) + 1
	currentTrackNo := uint(songsLength)
	var newPlaylistSongs []model.PlaylistSong
	for _, album := range albums {
		var currentSongIDs []uuid.UUID

		for _, song := range album.Songs {
			if slices.ContainsFunc(playlistSongs, func(p model.PlaylistSong) bool {
				return p.SongID == song.ID
			}) {
				currentSongIDs = append(currentSongIDs, song.ID)
				if request.ForceAdd != nil && !(*request.ForceAdd) {
					continue
				}
			}

			playlistSong := model.PlaylistSong{
				ID:          uuid.New(),
				PlaylistID:  request.ID,
				SongID:      song.ID,
				SongTrackNo: currentTrackNo,
			}
			newPlaylistSongs = append(newPlaylistSongs, playlistSong)
			currentTrackNo++
		}

		if len(currentSongIDs) == len(album.Songs) {
			duplicateAlbumIDs = append(duplicateAlbumIDs, album.ID)
		}
		duplicateSongIDs = append(duplicateSongIDs, currentSongIDs...)
	}

	if len(duplicateSongIDs) > 0 && request.ForceAdd == nil {
		return &responses.AddAlbumsToPlaylistResponse{
			Success:           false,
			DuplicateAlbumIDs: duplicateAlbumIDs,
			DuplicateSongIDs:  duplicateSongIDs,
		}, nil
	}

	err = a.repository.AddSongs(&newPlaylistSongs)
	if err != nil {
		return nil, wrapper.InternalServerError(err)
	}

	var addedSongIDs []uuid.UUID
	for _, ps := range newPlaylistSongs {
		addedSongIDs = append(addedSongIDs, ps.SongID)
	}

	return &responses.AddAlbumsToPlaylistResponse{
		Success:           true,
		DuplicateAlbumIDs: duplicateAlbumIDs,
		DuplicateSongIDs:  duplicateSongIDs,
		AddedSongIDs:      addedSongIDs,
	}, nil
}
