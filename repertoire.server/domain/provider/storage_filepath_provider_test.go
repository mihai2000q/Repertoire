package provider

import (
	"mime/multipart"
	"repertoire/server/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStorageFilePathProvider_GetAlbumImagePath_ShouldReturnAlbumImagePath(t *testing.T) {
	// given
	_uut := new(storageFilePathProvider)

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
	expectedImagePath := album.UserID.String() + "/albums/" + album.ID.String() + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetArtistImagePath_ShouldReturnArtistImagePath(t *testing.T) {
	// given
	_uut := new(storageFilePathProvider)

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
	expectedImagePath := artist.UserID.String() + "/artists/" + artist.ID.String() + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}

func TestStorageFilePathProvider_GetSongImagePath_ShouldReturnSongImagePath(t *testing.T) {
	// given
	_uut := new(storageFilePathProvider)

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
	expectedImagePath := song.UserID.String() + "/songs/" + song.ID.String() + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}
