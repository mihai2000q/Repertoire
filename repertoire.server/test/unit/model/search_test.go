package model

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
	"repertoire/server/model"
	"testing"
	"time"
)

func TestArtistSearch_ToSearch_WhenValid_ShouldReturnCorrectMapping(t *testing.T) {
	// given
	artist := model.Artist{
		ID:        uuid.New(),
		Name:      "Steve",
		ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    uuid.New(),
	}

	// when
	result := artist.ToSearch()

	// then
	assert.Equal(t, artist.ImageURL.StripURL(), result.ImageUrl)
	assert.Equal(t, artist.Name, result.Name)

	assert.Equal(t, "artist-"+artist.ID.String(), result.ID)
	assert.Equal(t, artist.CreatedAt, result.CreatedAt)
	assert.Equal(t, artist.UpdatedAt, result.UpdatedAt)
	assert.Equal(t, enums.Artist, result.Type)
	assert.Equal(t, artist.UserID, result.UserID)
}

func TestAlbumSearch_ToSearch_WhenValid_ShouldReturnCorrectMapping(t *testing.T) {
	tests := []struct {
		name  string
		album model.Album
	}{
		{
			"Without artist",
			model.Album{
				ID:        uuid.New(),
				Title:     "Some Album",
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			},
		},
		{
			"Full",
			model.Album{
				ID:          uuid.New(),
				Title:       "Some Album",
				ImageURL:    &[]internal.FilePath{"some_file_path"}[0],
				ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
				Artist: &model.Artist{
					ID:        uuid.New(),
					Name:      "Some Artist",
					ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
					UpdatedAt: time.Now().UTC(),
				},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given - parameterized
			// when
			result := tt.album.ToSearch()

			// then
			assert.Equal(t, tt.album.ImageURL.StripURL(), result.ImageUrl)
			assert.Equal(t, tt.album.Title, result.Title)
			if tt.album.ReleaseDate != nil {
				assert.NotNil(t, result.ReleaseDate)
				assert.Equal(t, (*time.Time)(tt.album.ReleaseDate).Format("2006-01-02"), *result.ReleaseDate)
			} else {
				assert.Nil(t, result.ReleaseDate)
			}

			if tt.album.Artist != nil {
				assert.Equal(t, tt.album.Artist.ID, result.Artist.ID)
				assert.Equal(t, tt.album.Artist.Name, result.Artist.Name)
				assert.Equal(t, tt.album.Artist.UpdatedAt, result.Artist.UpdatedAt)
				assert.Equal(t, tt.album.Artist.ImageURL.StripURL(), result.Artist.ImageUrl)
			} else {
				assert.Nil(t, result.Artist)
			}

			assert.Equal(t, "album-"+tt.album.ID.String(), result.ID)
			assert.Equal(t, tt.album.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.album.UpdatedAt, result.UpdatedAt)
			assert.Equal(t, enums.Album, result.Type)
			assert.Equal(t, tt.album.UserID, result.UserID)
		})
	}
}

func TestAlbumSearch_ToAlbumSearch_WhenValid_ShouldReturnCorrectMappingForArtist(t *testing.T) {
	// given
	artist := model.Artist{
		ID:        uuid.New(),
		Name:      "Steve",
		ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
		UpdatedAt: time.Now().UTC(),
	}

	// when
	result := artist.ToAlbumSearch()

	// then
	assert.Equal(t, artist.ImageURL.StripURL(), result.ImageUrl)
	assert.Equal(t, artist.Name, result.Name)
	assert.Equal(t, artist.ID, result.ID)
	assert.Equal(t, artist.UpdatedAt, result.UpdatedAt)
}

func TestSongSearch_ToSearch_WhenValid_ShouldReturnCorrectMapping(t *testing.T) {
	tests := []struct {
		name string
		song model.Song
	}{
		{
			"Only Mandatory",
			model.Song{
				ID:        uuid.New(),
				Title:     "Some Album",
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			},
		},
		{
			"Full",
			model.Song{
				ID:          uuid.New(),
				Title:       "Some Song",
				ImageURL:    &[]internal.FilePath{"some_file_path"}[0],
				ReleaseDate: &[]internal.Date{internal.Date(time.Now())}[0],
				Artist: &model.Artist{
					ID:        uuid.New(),
					Name:      "Some Artist",
					ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
					UpdatedAt: time.Now().UTC(),
				},
				Album: &model.Album{
					ID:        uuid.New(),
					Title:     "Some Album",
					ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
					UpdatedAt: time.Now().UTC(),
				},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
				UserID:    uuid.New(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given - parameterized
			// when
			result := tt.song.ToSearch()

			// then
			assert.Equal(t, tt.song.ImageURL.StripURL(), result.ImageUrl)
			assert.Equal(t, tt.song.Title, result.Title)
			if tt.song.ReleaseDate != nil {
				assert.NotNil(t, result.ReleaseDate)
				assert.Equal(t, (*time.Time)(tt.song.ReleaseDate).Format("2006-01-02"), *result.ReleaseDate)
			} else {
				assert.Nil(t, result.ReleaseDate)
			}

			if tt.song.Artist != nil {
				assert.Equal(t, tt.song.Artist.ID, result.Artist.ID)
				assert.Equal(t, tt.song.Artist.Name, result.Artist.Name)
				assert.Equal(t, tt.song.Artist.UpdatedAt, result.Artist.UpdatedAt)
				assert.Equal(t, tt.song.Artist.ImageURL.StripURL(), result.Artist.ImageUrl)
			} else {
				assert.Nil(t, result.Artist)
			}

			if tt.song.Album != nil {
				assert.Equal(t, tt.song.Album.ID, result.Album.ID)
				assert.Equal(t, tt.song.Album.Title, result.Album.Title)
				if tt.song.Album.ReleaseDate != nil {
					assert.NotNil(t, result.Album.ReleaseDate)
					assert.Equal(
						t,
						(*time.Time)(tt.song.Album.ReleaseDate).Format("2006-01-02"),
						*result.Album.ReleaseDate,
					)
				} else {
					assert.Nil(t, result.Album.ReleaseDate)
				}
				assert.Equal(t, tt.song.Album.UpdatedAt, result.Album.UpdatedAt)
				assert.Equal(t, tt.song.Album.ImageURL.StripURL(), result.Album.ImageUrl)
			} else {
				assert.Nil(t, result.Album)
			}

			assert.Equal(t, "song-"+tt.song.ID.String(), result.ID)
			assert.Equal(t, tt.song.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.song.UpdatedAt, result.UpdatedAt)
			assert.Equal(t, enums.Song, result.Type)
			assert.Equal(t, tt.song.UserID, result.UserID)
		})
	}
}

func TestSongSearch_ToSongSearch_WhenValid_ShouldReturnCorrectMappingForArtist(t *testing.T) {
	// given
	artist := model.Artist{
		ID:        uuid.New(),
		Name:      "Steve",
		ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
		UpdatedAt: time.Now().UTC(),
	}

	// when
	result := artist.ToSongSearch()

	// then
	assert.Equal(t, artist.ImageURL.StripURL(), result.ImageUrl)
	assert.Equal(t, artist.Name, result.Name)
	assert.Equal(t, artist.ID, result.ID)
	assert.Equal(t, artist.UpdatedAt, result.UpdatedAt)
}

func TestSongSearch_ToSongSearch_WhenValid_ShouldReturnCorrectMappingForAlbum(t *testing.T) {
	// given
	album := model.Album{
		ID:        uuid.New(),
		Title:     "Some Title",
		ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
		UpdatedAt: time.Now().UTC(),
	}

	// when
	result := album.ToSongSearch()

	// then
	assert.Equal(t, album.ImageURL.StripURL(), result.ImageUrl)
	assert.Equal(t, album.Title, result.Title)
	assert.Equal(t, album.ID, result.ID)
	assert.Equal(t, album.UpdatedAt, result.UpdatedAt)
}

func TestPlaylistSearch_ToSearch_WhenValid_ShouldReturnCorrectMapping(t *testing.T) {
	// given
	playlist := model.Playlist{
		ID:        uuid.New(),
		Title:     "Best Songs",
		ImageURL:  &[]internal.FilePath{"some_file_path"}[0],
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    uuid.New(),
	}

	// when
	result := playlist.ToSearch()

	// then
	assert.Equal(t, playlist.ImageURL.StripURL(), result.ImageUrl)
	assert.Equal(t, playlist.Title, result.Title)

	assert.Equal(t, "playlist-"+playlist.ID.String(), result.ID)
	assert.Equal(t, playlist.CreatedAt, result.CreatedAt)
	assert.Equal(t, playlist.UpdatedAt, result.UpdatedAt)
	assert.Equal(t, enums.Playlist, result.Type)
	assert.Equal(t, playlist.UserID, result.UserID)
}
