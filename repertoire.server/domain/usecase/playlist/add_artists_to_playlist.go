package playlist

import (
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
)

type AddArtistsToPlaylist struct {
	repository       repository.PlaylistRepository
	artistRepository repository.ArtistRepository
}

func NewAddArtistsToPlaylist(
	repository repository.PlaylistRepository,
	artistRepository repository.ArtistRepository,
) AddArtistsToPlaylist {
	return AddArtistsToPlaylist{
		repository:       repository,
		artistRepository: artistRepository,
	}
}

func (a AddArtistsToPlaylist) Handle(request requests.AddArtistsToPlaylistRequest) *wrapper.ErrorCode {
	var oldPlaylistSongs []model.PlaylistSong
	err := a.repository.GetPlaylistSongs(&oldPlaylistSongs, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	var artists []model.Artist
	err = a.artistRepository.GetAllByIDsWithSongs(&artists, request.ArtistIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	songsLength := len(oldPlaylistSongs) + 1
	currentTrackNo := uint(songsLength)
	var playlistSongs []model.PlaylistSong
	for _, artist := range artists {
		for _, song := range artist.Songs {
			if slices.ContainsFunc(oldPlaylistSongs, func(p model.PlaylistSong) bool {
				return p.SongID == song.ID
			}) {
				continue
			}

			playlistSong := model.PlaylistSong{
				ID:          uuid.New(),
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
