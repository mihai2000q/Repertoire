package provider

import (
	"mime/multipart"
	"path/filepath"
	"repertoire/server/model"
	"strings"
	"time"

	"github.com/google/uuid"
)

type StorageFilePathProvider interface {
	GetUserProfilePicturePath(file *multipart.FileHeader, user model.User) string
	GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string
	GetArtistImagePath(file *multipart.FileHeader, artist model.Artist) string
	GetBandMemberImagePath(file *multipart.FileHeader, artist model.BandMember) string
	GetPlaylistImagePath(file *multipart.FileHeader, playlist model.Playlist) string
	GetSongImagePath(file *multipart.FileHeader, song model.Song) string

	GetUserDirectoryPath(id uuid.UUID) string
	GetAlbumDirectoryPath(album model.Album) string
	GetArtistDirectoryPath(artist model.Artist) string
	GetPlaylistDirectoryPath(playlist model.Playlist) string
	GetSongDirectoryPath(song model.Song) string
}

type storageFilePathProvider struct {
}

func NewStorageFilePathProvider() StorageFilePathProvider {
	return new(storageFilePathProvider)
}

var albumRootDirectory = "albums"
var artistRootDirectory = "artists"
var bandMemberRootDirectory = "members"
var songRootDirectory = "songs"
var playlistRootDirectory = "playlists"

func (s storageFilePathProvider) GetUserProfilePicturePath(file *multipart.FileHeader, user model.User) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(user.ID.String()).
		WithFile(s.getTimeFormat(user.UpdatedAt) + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetAlbumImagePath(file *multipart.FileHeader, album model.Album) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(album.UserID.String()).
		WithDirectory(albumRootDirectory).
		WithDirectory(album.ID.String()).
		WithFile(s.getTimeFormat(album.UpdatedAt) + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetArtistImagePath(file *multipart.FileHeader, artist model.Artist) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.getArtistDirectory(artist).
		WithFile(s.getTimeFormat(artist.UpdatedAt) + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetBandMemberImagePath(file *multipart.FileHeader, member model.BandMember) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.getArtistDirectory(member.Artist).
		WithDirectory(bandMemberRootDirectory).
		WithDirectory(member.ID.String()).
		WithFile(s.getTimeFormat(member.UpdatedAt) + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetPlaylistImagePath(file *multipart.FileHeader, playlist model.Playlist) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(playlist.UserID.String()).
		WithDirectory(playlistRootDirectory).
		WithDirectory(playlist.ID.String()).
		WithFile(s.getTimeFormat(playlist.UpdatedAt) + fileExtension).
		BuildFilePath()
}

func (s storageFilePathProvider) GetSongImagePath(file *multipart.FileHeader, song model.Song) string {
	fileExtension := filepath.Ext(file.Filename)
	return s.builder().
		WithDirectory(song.UserID.String()).
		WithDirectory(songRootDirectory).
		WithDirectory(song.ID.String()).
		WithFile(s.getTimeFormat(song.UpdatedAt) + fileExtension).
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

func (s storageFilePathProvider) getTimeFormat(time time.Time) string {
	return time.Format("2006_01_02T15_04_05")
}

func (s storageFilePathProvider) getArtistDirectory(artist model.Artist) directoryPathBuilder {
	return s.builder().
		WithDirectory(artist.UserID.String()).
		WithDirectory(artistRootDirectory).
		WithDirectory(artist.ID.String())
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
