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

type AddArtistsToPlaylist struct {
	playlistRepository repository.PlaylistRepository
	artistRepository   repository.ArtistRepository
}

func NewAddArtistsToPlaylist(
	playlistRepository repository.PlaylistRepository,
	artistRepository repository.ArtistRepository,
) AddArtistsToPlaylist {
	return AddArtistsToPlaylist{
		playlistRepository: playlistRepository,
		artistRepository:   artistRepository,
	}
}

func (a AddArtistsToPlaylist) Handle(request requests.AddArtistsToPlaylistRequest) (*responses.AddArtistsToPlaylistResponse, *httperror.ErrorCode) {
	var playlistSongs []model.PlaylistSong
	if err := a.playlistRepository.GetPlaylistSongs(&playlistSongs, request.ID); err != nil {
		return nil, httperror.DatabaseError(err)
	}

	var artists []model.Artist
	if err := a.artistRepository.GetAllByIDsWithSongs(&artists, request.ArtistIDs); err != nil {
		return nil, httperror.DatabaseError(err)
	}
	if len(artists) != len(request.ArtistIDs) {
		return nil, httperror.NotFoundError(errors.New("artists not found"))
	}

	var duplicateSongIDs []uuid.UUID
	var duplicateArtistIDs []uuid.UUID

	songsLength := len(playlistSongs) + 1
	currentTrackNo := uint(songsLength)
	var newPlaylistSongs []model.PlaylistSong
	for _, artist := range artists {
		var currentSongIDs []uuid.UUID

		for _, song := range artist.Songs {
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

		if len(currentSongIDs) == len(artist.Songs) {
			duplicateArtistIDs = append(duplicateArtistIDs, artist.ID)
		}
		duplicateSongIDs = append(duplicateSongIDs, currentSongIDs...)
	}

	if len(duplicateSongIDs) > 0 && request.ForceAdd == nil {
		return &responses.AddArtistsToPlaylistResponse{
			Success:            false,
			DuplicateArtistIDs: duplicateArtistIDs,
			DuplicateSongIDs:   duplicateSongIDs,
		}, nil
	}

	if err := a.playlistRepository.AddSongs(&newPlaylistSongs); err != nil {
		return nil, httperror.DatabaseError(err)
	}

	var addedSongIDs []uuid.UUID
	for _, ps := range newPlaylistSongs {
		addedSongIDs = append(addedSongIDs, ps.SongID)
	}

	return &responses.AddArtistsToPlaylistResponse{
		Success:            true,
		DuplicateArtistIDs: duplicateArtistIDs,
		DuplicateSongIDs:   duplicateSongIDs,
		AddedSongIDs:       addedSongIDs,
	}, nil
}
