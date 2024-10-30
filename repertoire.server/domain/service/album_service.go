package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"

	"github.com/google/uuid"
)

type AlbumService interface {
	Get(id uuid.UUID) (model.Album, *wrapper.ErrorCode)
	GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.Album], *wrapper.ErrorCode)
	Create(request requests.CreateAlbumRequest, token string) *wrapper.ErrorCode
	Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
}

type albumService struct {
	getAlbum     album.GetAlbum
	getAllAlbums album.GetAllAlbums
	createAlbum  album.CreateAlbum
	updateAlbum  album.UpdateAlbum
	deleteAlbum  album.DeleteAlbum
}

func NewAlbumService(
	getAlbum album.GetAlbum,
	getAllAlbums album.GetAllAlbums,
	createAlbum album.CreateAlbum,
	updateAlbum album.UpdateAlbum,
	deleteAlbum album.DeleteAlbum,
) AlbumService {
	return &albumService{
		getAlbum:     getAlbum,
		getAllAlbums: getAllAlbums,
		createAlbum:  createAlbum,
		updateAlbum:  updateAlbum,
		deleteAlbum:  deleteAlbum,
	}
}

func (a *albumService) Get(id uuid.UUID) (model.Album, *wrapper.ErrorCode) {
	return a.getAlbum.Handle(id)
}

func (a *albumService) GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.Album], *wrapper.ErrorCode) {
	return a.getAllAlbums.Handle(request, token)
}

func (a *albumService) Create(request requests.CreateAlbumRequest, token string) *wrapper.ErrorCode {
	return a.createAlbum.Handle(request, token)
}

func (a *albumService) Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode {
	return a.updateAlbum.Handle(request)
}

func (a *albumService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteAlbum.Handle(id)
}
