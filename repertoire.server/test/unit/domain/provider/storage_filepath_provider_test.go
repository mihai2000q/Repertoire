package provider

import (
	"mime/multipart"
	"repertoire/server/domain/provider"
	"repertoire/server/internal"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStorageFilePathProvider_GetUserProfilePicturePath_ShouldReturnPlaylistImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	user := model.User{ID: uuid.New()}

	// when
	imagePath := _uut.GetUserProfilePicturePath(file, user)

	// then
	expectedImagePath := user.ID.String() + "/profile_pic" + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetAlbumImagePath_ShouldReturnAlbumImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	album := model.Album{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	imagePath := _uut.GetAlbumImagePath(file, album)

	// then
	expectedImagePath := album.UserID.String() + "/albums/" + album.ID.String() + "/image" + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetArtistImagePath_ShouldReturnArtistImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	artist := model.Artist{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	imagePath := _uut.GetArtistImagePath(file, artist)

	// then
	expectedImagePath := artist.UserID.String() + "/artists/" + artist.ID.String() + "/image" + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetPlaylistImagePath_ShouldReturnPlaylistImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	playlist := model.Playlist{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	imagePath := _uut.GetPlaylistImagePath(file, playlist)

	// then
	expectedImagePath := playlist.UserID.String() + "/playlists/" + playlist.ID.String() + "/image" + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetSongImagePath_ShouldReturnSongImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	song := model.Song{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	imagePath := _uut.GetSongImagePath(file, song)

	// then
	expectedImagePath := song.UserID.String() + "/songs/" + song.ID.String() + "/image" + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetUserDirectoryPath_ShouldReturnUserDirectoryPath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	userID := uuid.New()

	// when
	directoryPath := _uut.GetUserDirectoryPath(userID)

	// then
	expectedDirectoryPath := userID.String() + "/"

	assert.Equal(t, expectedDirectoryPath, directoryPath)
}

func TestStorageFilePathProvider_GetAlbumDirectoryPath_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	album := model.Album{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	directoryPath := _uut.GetAlbumDirectoryPath(album)

	// then
	expectedDirectoryPath := album.UserID.String() + "/albums/" + album.ID.String() + "/"

	assert.Equal(t, expectedDirectoryPath, directoryPath)
}

func TestStorageFilePathProvider_GetArtistDirectoryPath_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	artist := model.Artist{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	directoryPath := _uut.GetArtistDirectoryPath(artist)

	// then
	expectedDirectoryPath := artist.UserID.String() + "/artists/" + artist.ID.String() + "/"

	assert.Equal(t, expectedDirectoryPath, directoryPath)
}

func TestStorageFilePathProvider_GetPlaylistDirectoryPath_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	playlist := model.Playlist{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	directoryPath := _uut.GetPlaylistDirectoryPath(playlist)

	// then
	expectedDirectoryPath := playlist.UserID.String() + "/playlists/" + playlist.ID.String() + "/"

	assert.Equal(t, expectedDirectoryPath, directoryPath)
}

func TestStorageFilePathProvider_GetSongDirectoryPath_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	song := model.Song{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}

	// when
	directoryPath := _uut.GetSongDirectoryPath(song)

	// then
	expectedDirectoryPath := song.UserID.String() + "/songs/" + song.ID.String() + "/"

	assert.Equal(t, expectedDirectoryPath, directoryPath)
}

func TestStorageFilePathProvider_HasAlbumFiles_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	tests := []struct {
		name  string
		album model.Album
		want  bool
	}{
		{
			"With Image",
			model.Album{
				ImageURL: &[]internal.FilePath{"somePath"}[0],
			},
			true,
		},
		{
			"Without Files",
			model.Album{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := provider.NewStorageFilePathProvider()

			// when
			result := _uut.HasAlbumFiles(tt.album)

			// then
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStorageFilePathProvider_HasArtistFiles_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	tests := []struct {
		name   string
		artist model.Artist
		want   bool
	}{
		{
			"With Image",
			model.Artist{
				ImageURL: &[]internal.FilePath{"somePath"}[0],
			},
			true,
		},
		{
			"Without Files",
			model.Artist{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := provider.NewStorageFilePathProvider()

			// when
			result := _uut.HasArtistFiles(tt.artist)

			// then
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStorageFilePathProvider_HasPlaylistFiles_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	tests := []struct {
		name     string
		playlist model.Playlist
		want     bool
	}{
		{
			"With Image",
			model.Playlist{
				ImageURL: &[]internal.FilePath{"somePath"}[0],
			},
			true,
		},
		{
			"Without Files",
			model.Playlist{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := provider.NewStorageFilePathProvider()

			// when
			result := _uut.HasPlaylistFiles(tt.playlist)

			// then
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStorageFilePathProvider_HasSongFiles_ShouldReturnAlbumDirectoryPath(t *testing.T) {
	tests := []struct {
		name string
		song model.Song
		want bool
	}{
		{
			"With Image",
			model.Song{
				ImageURL: &[]internal.FilePath{"somePath"}[0],
			},
			true,
		},
		{
			"Without Files",
			model.Song{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := provider.NewStorageFilePathProvider()

			// when
			result := _uut.HasSongFiles(tt.song)

			// then
			assert.Equal(t, tt.want, result)
		})
	}
}
