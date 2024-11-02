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
	songImagePath, songImageUrl := _uut.GetSongImagePathAndURL(file, songID)

	// then
	expectedImagePath := "songs/" + songID.String() + fileExtension
	expectedImageUrl := env.StorageUrl + "/files/" + expectedImagePath

	assert.Equal(t, expectedImageUrl, songImageUrl)
	assert.Equal(t, expectedImagePath, songImagePath)
}
