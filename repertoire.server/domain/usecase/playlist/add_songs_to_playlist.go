package playlist

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/api/responses"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
)

type AddSongsToPlaylist struct {
	repository repository.PlaylistRepository
}

func NewAddSongsToPlaylist(repository repository.PlaylistRepository) AddSongsToPlaylist {
	return AddSongsToPlaylist{repository: repository}
}

func (a AddSongsToPlaylist) Handle(
	request requests.AddSongsToPlaylistRequest,
) (*responses.AddSongsToPlaylistResponse, *wrapper.ErrorCode) {
	var playlistSongs []model.PlaylistSong
	err := a.repository.GetPlaylistSongs(&playlistSongs, request.ID)
	if err != nil {
		return nil, wrapper.InternalServerError(err)
	}

	var duplicates []uuid.UUID
	var newPlaylistSongs []model.PlaylistSong
	currentTrackNo := uint(len(playlistSongs)) + 1
	for _, songID := range request.SongIDs {
		if slices.ContainsFunc(playlistSongs, func(p model.PlaylistSong) bool {
			return p.SongID == songID
		}) {
			duplicates = append(duplicates, songID)
			if request.ForceAdd != nil && !(*request.ForceAdd) {
				continue
			}
		}

		playlistSong := model.PlaylistSong{
			ID:          uuid.New(),
			PlaylistID:  request.ID,
			SongID:      songID,
			SongTrackNo: currentTrackNo,
		}
		newPlaylistSongs = append(newPlaylistSongs, playlistSong)
		currentTrackNo++
	}

	if len(duplicates) == 0 && request.ForceAdd != nil {
		return nil, wrapper.BadRequestError(errors.New("force adding when there are no duplicates"))
	}
	if len(duplicates) > 0 && request.ForceAdd == nil {
		return &responses.AddSongsToPlaylistResponse{Success: false, Duplicates: duplicates}, nil
	}

	err = a.repository.AddSongs(&newPlaylistSongs)
	if err != nil {
		return nil, wrapper.InternalServerError(err)
	}

	var addedSongs []uuid.UUID
	for _, ps := range newPlaylistSongs {
		addedSongs = append(addedSongs, ps.SongID)
	}

	return &responses.AddSongsToPlaylistResponse{
		Success:    true,
		Duplicates: duplicates,
		Added:      addedSongs,
	}, nil
}
