package provider

import (
	"mime/multipart"
	"path/filepath"
	"repertoire/server/model"
)

type StorageFilePathProvider interface {
	GetUserProfilePicturePath(file *multipart.FileHeader, user model.User) string
	GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string
	GetArtistImagePath(file *multipart.FileHeader, artist model.Artist) string
	GetPlaylistImagePath(file *multipart.FileHeader, playlist model.Playlist) string
	GetSongImagePath(file *multipart.FileHeader, song model.Song) string
}

type storageFilePathProvider struct {
}

func NewStorageFilePathProvider() StorageFilePathProvider {
	return new(storageFilePathProvider)
}

func (s storageFilePathProvider) GetUserProfilePicturePath(file *multipart.FileHeader, user model.User) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := user.ID.String() + "/profile_pic" + fileExtension
	return filePath
}

func (s storageFilePathProvider) GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := album.UserID.String() + "/albums/" + album.ID.String() + fileExtension
	return filePath
}

func (s storageFilePathProvider) GetArtistImagePath(file *multipart.FileHeader, artist model.Artist) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := artist.UserID.String() + "/artists/" + artist.ID.String() + fileExtension
	return filePath
}

func (s storageFilePathProvider) GetPlaylistImagePath(file *multipart.FileHeader, playlist model.Playlist) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := playlist.UserID.String() + "/playlists/" + playlist.ID.String() + fileExtension
	return filePath
}

func (s storageFilePathProvider) GetSongImagePath(file *multipart.FileHeader, song model.Song) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := song.UserID.String() + "/songs/" + song.ID.String() + fileExtension
	return filePath
}
