package provider

import (
	"mime/multipart"
	"path/filepath"
	"repertoire/server/model"
)

type StorageFilePathProvider interface {
	GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string
	GetSongImagePath(file *multipart.FileHeader, song model.Song) string
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

func (s storageFilePathProvider) GetSongImagePath(file *multipart.FileHeader, song model.Song) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := song.UserID.String() + "/songs/" + song.ID.String() + fileExtension
	return filePath
}
