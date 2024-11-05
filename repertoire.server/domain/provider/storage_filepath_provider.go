package provider

import (
	"mime/multipart"
	"path/filepath"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type StorageFilePathProvider interface {
	GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string
	GetSongImagePath(file *multipart.FileHeader, songID uuid.UUID) string
}

type storageFilePathProvider struct {
}

func NewStorageFilePathProvider() StorageFilePathProvider {
	return new(storageFilePathProvider)
}

func (s storageFilePathProvider) GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := album.UserID.String() + "/albums/" + album.ID.String() + fileExtension
	return filePath
}

func (s storageFilePathProvider) GetSongImagePath(file *multipart.FileHeader, songID uuid.UUID) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := "songs/" + songID.String() + fileExtension
	return filePath
}
