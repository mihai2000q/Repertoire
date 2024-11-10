package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type AlbumService interface {
	AddSong(requests.AddSongToAlbumRequest) *wrapper.ErrorCode
	Create(request requests.CreateAlbumRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	Delete(id uuid.UUID) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	Get(id uuid.UUID) (model.Album, *wrapper.ErrorCode)
	GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.Album], *wrapper.ErrorCode)
	MoveSong(request requests.MoveSongFromAlbumRequest) *wrapper.ErrorCode
	RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode
}

type albumService struct {
	addSongToAlbum       album.AddSongToAlbum
	createAlbum          album.CreateAlbum
	deleteAlbum          album.DeleteAlbum
	deleteImageFromAlbum album.DeleteImageFromAlbum
	getAlbum             album.GetAlbum
	getAllAlbums         album.GetAllAlbums
	moveSongFromAlbum    album.MoveSongFromAlbum
	removeSongFromAlbum  album.RemoveSongFromAlbum
	saveImageToAlbum     album.SaveImageToAlbum
	updateAlbum          album.UpdateAlbum
}

func NewAlbumService(
	addSongToAlbum album.AddSongToAlbum,
	createAlbum album.CreateAlbum,
	deleteAlbum album.DeleteAlbum,
	deleteImageFromAlbum album.DeleteImageFromAlbum,
	getAlbum album.GetAlbum,
	getAllAlbums album.GetAllAlbums,
	moveSongFromAlbum album.MoveSongFromAlbum,
	removeSongFromAlbum album.RemoveSongFromAlbum,
	saveImageToAlbum album.SaveImageToAlbum,
	updateAlbum album.UpdateAlbum,
) AlbumService {
	return &albumService{
		addSongToAlbum:       addSongToAlbum,
		createAlbum:          createAlbum,
		deleteAlbum:          deleteAlbum,
		deleteImageFromAlbum: deleteImageFromAlbum,
		getAlbum:             getAlbum,
		getAllAlbums:         getAllAlbums,
		moveSongFromAlbum:    moveSongFromAlbum,
		removeSongFromAlbum:  removeSongFromAlbum,
		saveImageToAlbum:     saveImageToAlbum,
		updateAlbum:          updateAlbum,
	}
}

func (a *albumService) AddSong(request requests.AddSongToAlbumRequest) *wrapper.ErrorCode {
	return a.addSongToAlbum.Handle(request)
}

func (a *albumService) Create(request requests.CreateAlbumRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	return a.createAlbum.Handle(request, token)
}

func (a *albumService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteAlbum.Handle(id)
}

func (a *albumService) DeleteImage(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteImageFromAlbum.Handle(id)
}

func (a *albumService) Get(id uuid.UUID) (model.Album, *wrapper.ErrorCode) {
	return a.getAlbum.Handle(id)
}

func (a *albumService) GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.Album], *wrapper.ErrorCode) {
	return a.getAllAlbums.Handle(request, token)
}

func (a *albumService) MoveSong(request requests.MoveSongFromAlbumRequest) *wrapper.ErrorCode {
	return a.moveSongFromAlbum.Handle(request)
}

func (a *albumService) RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return a.removeSongFromAlbum.Handle(id, songID)
}

func (a *albumService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return a.saveImageToAlbum.Handle(file, id)
}

func (a *albumService) Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode {
	return a.updateAlbum.Handle(request)
}
