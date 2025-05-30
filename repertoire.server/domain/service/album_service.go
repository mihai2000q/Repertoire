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
	AddSongs(requests.AddSongsToAlbumRequest) *wrapper.ErrorCode
	Create(request requests.CreateAlbumRequest, token string) (uuid.UUID, *wrapper.ErrorCode)
	Delete(request requests.DeleteAlbumRequest) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	Get(request requests.GetAlbumRequest) (model.Album, *wrapper.ErrorCode)
	GetFiltersMetadata(
		request requests.GetAlbumFiltersMetadataRequest,
		token string,
	) (model.AlbumFiltersMetadata, *wrapper.ErrorCode)
	GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.EnhancedAlbum], *wrapper.ErrorCode)
	MoveSong(request requests.MoveSongFromAlbumRequest) *wrapper.ErrorCode
	RemoveSongs(request requests.RemoveSongsFromAlbumRequest) *wrapper.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode
}

type albumService struct {
	addSongsToAlbum         album.AddSongsToAlbum
	createAlbum             album.CreateAlbum
	deleteAlbum             album.DeleteAlbum
	deleteImageFromAlbum    album.DeleteImageFromAlbum
	getAlbum                album.GetAlbum
	getAlbumFiltersMetadata album.GetAlbumFiltersMetadata
	getAllAlbums            album.GetAllAlbums
	moveSongFromAlbum       album.MoveSongFromAlbum
	removeSongsFromAlbum    album.RemoveSongsFromAlbum
	saveImageToAlbum        album.SaveImageToAlbum
	updateAlbum             album.UpdateAlbum
}

func NewAlbumService(
	addSongsToAlbum album.AddSongsToAlbum,
	createAlbum album.CreateAlbum,
	deleteAlbum album.DeleteAlbum,
	deleteImageFromAlbum album.DeleteImageFromAlbum,
	getAlbum album.GetAlbum,
	getAlbumFiltersMetadata album.GetAlbumFiltersMetadata,
	getAllAlbums album.GetAllAlbums,
	moveSongFromAlbum album.MoveSongFromAlbum,
	removeSongsFromAlbum album.RemoveSongsFromAlbum,
	saveImageToAlbum album.SaveImageToAlbum,
	updateAlbum album.UpdateAlbum,
) AlbumService {
	return &albumService{
		addSongsToAlbum:         addSongsToAlbum,
		createAlbum:             createAlbum,
		deleteAlbum:             deleteAlbum,
		deleteImageFromAlbum:    deleteImageFromAlbum,
		getAlbum:                getAlbum,
		getAlbumFiltersMetadata: getAlbumFiltersMetadata,
		getAllAlbums:            getAllAlbums,
		moveSongFromAlbum:       moveSongFromAlbum,
		removeSongsFromAlbum:    removeSongsFromAlbum,
		saveImageToAlbum:        saveImageToAlbum,
		updateAlbum:             updateAlbum,
	}
}

func (a *albumService) AddSongs(request requests.AddSongsToAlbumRequest) *wrapper.ErrorCode {
	return a.addSongsToAlbum.Handle(request)
}

func (a *albumService) Create(request requests.CreateAlbumRequest, token string) (uuid.UUID, *wrapper.ErrorCode) {
	return a.createAlbum.Handle(request, token)
}

func (a *albumService) Delete(request requests.DeleteAlbumRequest) *wrapper.ErrorCode {
	return a.deleteAlbum.Handle(request)
}

func (a *albumService) DeleteImage(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteImageFromAlbum.Handle(id)
}

func (a *albumService) Get(request requests.GetAlbumRequest) (model.Album, *wrapper.ErrorCode) {
	return a.getAlbum.Handle(request)
}

func (a *albumService) GetFiltersMetadata(
	request requests.GetAlbumFiltersMetadataRequest,
	token string,
) (model.AlbumFiltersMetadata, *wrapper.ErrorCode) {
	return a.getAlbumFiltersMetadata.Handle(request, token)
}

func (a *albumService) GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.EnhancedAlbum], *wrapper.ErrorCode) {
	return a.getAllAlbums.Handle(request, token)
}

func (a *albumService) MoveSong(request requests.MoveSongFromAlbumRequest) *wrapper.ErrorCode {
	return a.moveSongFromAlbum.Handle(request)
}

func (a *albumService) RemoveSongs(request requests.RemoveSongsFromAlbumRequest) *wrapper.ErrorCode {
	return a.removeSongsFromAlbum.Handle(request)
}

func (a *albumService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return a.saveImageToAlbum.Handle(file, id)
}

func (a *albumService) Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode {
	return a.updateAlbum.Handle(request)
}
