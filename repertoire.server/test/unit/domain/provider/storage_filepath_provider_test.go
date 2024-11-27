package provider

import (
	"mime/multipart"
	provider2 "repertoire/server/domain/provider"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStorageFilePathProvider_GetUserProfilePicturePath_ShouldReturnPlaylistImagePath(t *testing.T) {
	// given
	_uut := new(provider2.storageFilePathProvider)

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
	_uut := new(provider2.storageFilePathProvider)

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
	_uut := new(provider2.storageFilePathProvider)

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
	_uut := new(provider2.storageFilePathProvider)

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
	_uut := new(provider2.storageFilePathProvider)

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
