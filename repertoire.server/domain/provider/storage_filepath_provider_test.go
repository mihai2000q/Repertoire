package provider

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"repertoire/server/internal"
	"testing"
)

func TestStorageFilePathProvider_GetSongImagePathAndURL_ShouldReturnImagePathAndUrl(t *testing.T) {
	// given
	env := internal.Env{StorageUrl: "storageUrl"}
	_uut := storageFilePathProvider{env}

	fileExtension := ".jpg"
	file := new(multipart.FileHeader)
	file.Filename = "something" + fileExtension
	songID := uuid.New()

	// when
	songImagePath := _uut.GetSongImagePath(file, songID)

	// then
	expectedImagePath := "songs/" + songID.String() + fileExtension

	assert.Equal(t, expectedImagePath, songImagePath)
}
