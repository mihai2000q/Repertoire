package provider

import (
	"mime/multipart"
	"repertoire/server/domain/provider"
	"repertoire/server/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStorageFilePathProvider_GetUserProfilePicturePath_ShouldReturnPlaylistImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	user := model.User{ID: uuid.New(), UpdatedAt: time.Now()}

	// when
	imagePath := _uut.GetUserProfilePicturePath(file, user)

	// then
	expectedImagePath := user.ID.String() + "/" +
		user.UpdatedAt.Format("2006-01-02T15:04:05") + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetAlbumImagePath_ShouldReturnAlbumImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	album := model.Album{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		UpdatedAt: time.Now(),
	}

	// when
	imagePath := _uut.GetAlbumImagePath(file, album)

	// then
	expectedImagePath := album.UserID.String() + "/albums/" + album.ID.String() + "/" +
		album.UpdatedAt.Format("2006-01-02T15:04:05") + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetArtistImagePath_ShouldReturnArtistImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	artist := model.Artist{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		UpdatedAt: time.Now(),
	}

	// when
	imagePath := _uut.GetArtistImagePath(file, artist)

	// then
	expectedImagePath := artist.UserID.String() + "/artists/" + artist.ID.String() + "/" +
		artist.UpdatedAt.Format("2006-01-02T15:04:05") + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetBandMemberImagePath_ShouldReturnBandMemberImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	artist := model.Artist{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}
	bandMember := model.BandMember{
		ID:        uuid.New(),
		Artist:    artist,
		UpdatedAt: time.Now(),
	}

	// when
	imagePath := _uut.GetBandMemberImagePath(file, bandMember)

	// then
	expectedImagePath := artist.UserID.String() + "/artists/" + artist.ID.String() +
		"/members/" + bandMember.ID.String() + "/" +
		bandMember.UpdatedAt.Format("2006-01-02T15:04:05") + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetPlaylistImagePath_ShouldReturnPlaylistImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	playlist := model.Playlist{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		UpdatedAt: time.Now(),
	}

	// when
	imagePath := _uut.GetPlaylistImagePath(file, playlist)

	// then
	expectedImagePath := playlist.UserID.String() + "/playlists/" + playlist.ID.String() + "/" +
		playlist.UpdatedAt.Format("2006-01-02T15:04:05") + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetSongImagePath_ShouldReturnSongImagePath(t *testing.T) {
	// given
	_uut := provider.NewStorageFilePathProvider()

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	song := model.Song{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		UpdatedAt: time.Now(),
	}

	// when
	imagePath := _uut.GetSongImagePath(file, song)

	// then
	expectedImagePath := song.UserID.String() + "/songs/" + song.ID.String() + "/" +
		song.UpdatedAt.Format("2006-01-02T15:04:05") + fileExtension

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
