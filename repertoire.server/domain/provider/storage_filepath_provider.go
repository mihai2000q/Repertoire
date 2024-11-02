package provider

import (
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
	"repertoire/server/internal"
)

type StorageFilePathProvider interface {
	GetSongImagePathAndURL(file *multipart.FileHeader, songID uuid.UUID) (string, string)
}

type storageFilePathProvider struct {
	env internal.Env
}

func NewStorageFilePathProvider(env internal.Env) StorageFilePathProvider {
	return storageFilePathProvider{
		env: env,
	}
}

func (s storageFilePathProvider) GetSongImagePathAndURL(file *multipart.FileHeader, songID uuid.UUID) (string, string) {
	fileExtension := filepath.Ext(file.Filename)
	filePath := "songs/" + songID.String() + fileExtension
	return filePath, s.getBaseURL() + filePath
}

func (s storageFilePathProvider) getBaseURL() string {
	return s.env.StorageUrl + "/files/"
}
