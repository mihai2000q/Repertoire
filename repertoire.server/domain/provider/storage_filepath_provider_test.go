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

func TestStorageFilePathProvider_GetSongImagePath_ShouldReturnSongImagePath(t *testing.T) {
	// given
	_uut := new(storageFilePathProvider)

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	songID := uuid.New()

	// when
	imagePath := _uut.GetSongImagePath(file, songID)

	// then
	expectedImagePath := "songs/" + songID.String() + fileExtension

	assert.Equal(t, expectedImagePath, imagePath)
}
