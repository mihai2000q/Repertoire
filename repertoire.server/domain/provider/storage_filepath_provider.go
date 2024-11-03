package provider

import (
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
	"repertoire/server/internal"
)

type StorageFilePathProvider interface {
	GetSongImagePath(file *multipart.FileHeader, songID uuid.UUID) string
}

type storageFilePathProvider struct {
	env internal.Env
}

func NewStorageFilePathProvider(env internal.Env) StorageFilePathProvider {
	return storageFilePathProvider{
		env: env,
	}
}

func (s storageFilePathProvider) GetSongImagePath(file *multipart.FileHeader, songID uuid.UUID) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := "songs/" + songID.String() + fileExtension
	return filePath
}
