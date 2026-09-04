package playlist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/api/responses"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
	"slices"

	"github.com/google/uuid"
)

type AddAlbumsToPlaylist struct {
	playlistRepository repository.PlaylistRepository
	albumRepository    repository.AlbumRepository
}

func NewAddAlbumsToPlaylist(
	playlistRepository repository.PlaylistRepository,
	albumRepository repository.AlbumRepository,
) AddAlbumsToPlaylist {
	return AddAlbumsToPlaylist{
		playlistRepository: playlistRepository,
		albumRepository:    albumRepository,
	}
}

func (a AddAlbumsToPlaylist) Handle(
	request requests.AddAlbumsToPlaylistRequest,
) (*responses.AddAlbumsToPlaylistResponse, *httperror.ErrorCode) {
	var playlistSongs []model.PlaylistSong
	if err := a.playlistRepository.GetPlaylistSongs(&playlistSongs, request.ID); err != nil {
		return nil, httperror.DatabaseError(err)
	}

	var albums []model.Album
	if err := a.albumRepository.GetAllByIDsWithSongs(&albums, request.AlbumIDs); err != nil {
		return nil, httperror.DatabaseError(err)
	}
	if len(albums) != len(request.AlbumIDs) {
		return nil, httperror.NotFoundError(errors.New("albums not found"))
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

	if err := a.playlistRepository.AddSongs(&newPlaylistSongs); err != nil {
		return nil, httperror.DatabaseError(err)
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
