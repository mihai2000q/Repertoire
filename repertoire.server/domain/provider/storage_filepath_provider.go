package provider

import (
	"github.com/google/uuid"
	"mime/multipart"
	"path/filepath"
	"repertoire/server/model"
	"strings"
)

type StorageFilePathProvider interface {
	GetUserProfilePicturePath(file *multipart.FileHeader, user model.User) string
	GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string
	GetArtistImagePath(file *multipart.FileHeader, artist model.Artist) string
	GetPlaylistImagePath(file *multipart.FileHeader, playlist model.Playlist) string
	GetSongImagePath(file *multipart.FileHeader, song model.Song) string

	GetUserDirectoryPath(id uuid.UUID) string
	GetAlbumDirectoryPath(album model.Album) string
	GetArtistDirectoryPath(artist model.Artist) string
	GetPlaylistDirectoryPath(playlist model.Playlist) string
	GetSongDirectoryPath(song model.Song) string

	HasAlbumFiles(album model.Album) bool
	HasArtistFiles(artist model.Artist) bool
	HasPlaylistFiles(playlist model.Playlist) bool
	HasSongFiles(song model.Song) bool
}

type storageFilePathProvider struct {
}

func NewStorageFilePathProvider() StorageFilePathProvider {
	return new(storageFilePathProvider)
}

var profilePicture = "profile_pic"
var image = "image"
var albumRootDirectory = "albums"
var artistRootDirectory = "artists"
var songRootDirectory = "songs"
var playlistRootDirectory = "playlists"

func (s storageFilePathProvider) GetUserProfilePicturePath(file *multipart.FileHeader, user model.User) string {
	fileExtension := filepath.Ext(file.Filename)
	filePath := user.ID.String() + "/" + profilePicture + fileExtension
	return filePath
}

func (s storageFilePathProvider) GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(album.UserID.String()).
		WithDirectory(albumRootDirectory).
		WithDirectory(album.ID.String()).
		WithFile(image + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetArtistImagePath(file *multipart.FileHeader, artist model.Artist) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(artist.UserID.String()).
		WithDirectory(artistRootDirectory).
		WithDirectory(artist.ID.String()).
		WithFile(image + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetPlaylistImagePath(file *multipart.FileHeader, playlist model.Playlist) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(playlist.UserID.String()).
		WithDirectory(playlistRootDirectory).
		WithDirectory(playlist.ID.String()).
		WithFile(image + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetSongImagePath(file *multipart.FileHeader, song model.Song) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(song.UserID.String()).
		WithDirectory(songRootDirectory).
		WithDirectory(song.ID.String()).
		WithFile(image + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetUserDirectoryPath(userID uuid.UUID) string {
	return s.builder().
		WithDirectory(userID.String()).
		BuildDirectoryPath()
}

func (s storageFilePathProvider) GetAlbumDirectoryPath(album model.Album) string {
	return s.builder().
		WithDirectory(album.UserID.String()).
		WithDirectory(albumRootDirectory).
		WithDirectory(album.ID.String()).
		BuildDirectoryPath()
}

func (s storageFilePathProvider) GetArtistDirectoryPath(artist model.Artist) string {
	return s.builder().
		WithDirectory(artist.UserID.String()).
		WithDirectory(artistRootDirectory).
		WithDirectory(artist.ID.String()).
		BuildDirectoryPath()
}

func (s storageFilePathProvider) GetPlaylistDirectoryPath(playlist model.Playlist) string {
	return s.builder().
		WithDirectory(playlist.UserID.String()).
		WithDirectory(playlistRootDirectory).
		WithDirectory(playlist.ID.String()).
		BuildDirectoryPath()
}

func (s storageFilePathProvider) GetSongDirectoryPath(song model.Song) string {
	return s.builder().
		WithDirectory(song.UserID.String()).
		WithDirectory(songRootDirectory).
		WithDirectory(song.ID.String()).
		BuildDirectoryPath()
}

func (s storageFilePathProvider) HasAlbumFiles(album model.Album) bool {
	return album.ImageURL != nil
}

func (s storageFilePathProvider) HasArtistFiles(artist model.Artist) bool {
	return artist.ImageURL != nil
}

func (s storageFilePathProvider) HasPlaylistFiles(playlist model.Playlist) bool {
	return playlist.ImageURL != nil
}

func (s storageFilePathProvider) HasSongFiles(song model.Song) bool {
	return song.ImageURL != nil
}

func (s storageFilePathProvider) builder() directoryPathBuilder {
	return pathBuilder{}
}

type directoryPathBuilder interface {
	WithDirectory(directory string) directoryPathBuilder
	WithFile(file string) filePathBuilder
	BuildDirectoryPath() string
}

type filePathBuilder interface {
	BuildFilePath() string
}

type pathBuilder struct {
	directories []string
	file        string
}

func (p pathBuilder) WithDirectory(directory string) directoryPathBuilder {
	p.directories = append(p.directories, directory)
	return p
}

func (p pathBuilder) WithFile(file string) filePathBuilder {
	p.file = file
	return p
}

func (p pathBuilder) BuildFilePath() string {
	return strings.Join(p.directories, "/") + "/" + p.file
}

func (p pathBuilder) BuildDirectoryPath() string {
	return strings.Join(p.directories, "/") + "/"
}
