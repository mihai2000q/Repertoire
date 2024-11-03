package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type AlbumService interface {
	Create(request requests.CreateAlbumRequest, token string) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
	Get(id uuid.UUID) (model.Album, *wrapper.ErrorCode)
	GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.Album], *wrapper.ErrorCode)
	RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode
}

type albumService struct {
	createAlbum         album.CreateAlbum
	deleteAlbum         album.DeleteAlbum
	getAlbum            album.GetAlbum
	getAllAlbums        album.GetAllAlbums
	removeSongFromAlbum album.RemoveSongFromAlbum
	updateAlbum         album.UpdateAlbum
}

func NewAlbumService(
	createAlbum album.CreateAlbum,
	deleteAlbum album.DeleteAlbum,
	getAlbum album.GetAlbum,
	getAllAlbums album.GetAllAlbums,
	removeSongFromAlbum album.RemoveSongFromAlbum,
	updateAlbum album.UpdateAlbum,
) AlbumService {
	return &albumService{
		createAlbum:         createAlbum,
		deleteAlbum:         deleteAlbum,
		getAlbum:            getAlbum,
		getAllAlbums:        getAllAlbums,
		removeSongFromAlbum: removeSongFromAlbum,
		updateAlbum:         updateAlbum,
	}
}

func (a *albumService) Create(request requests.CreateAlbumRequest, token string) *wrapper.ErrorCode {
	return a.createAlbum.Handle(request, token)
}

func (a *albumService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteAlbum.Handle(id)
}

func (a *albumService) Get(id uuid.UUID) (model.Album, *wrapper.ErrorCode) {
	return a.getAlbum.Handle(id)
}

func (a *albumService) GetAll(request requests.GetAlbumsRequest, token string) (wrapper.WithTotalCount[model.Album], *wrapper.ErrorCode) {
	return a.getAllAlbums.Handle(request, token)
}

func (a *albumService) RemoveSong(id uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	return a.removeSongFromAlbum.Handle(id, songID)
}

func (a *albumService) Update(request requests.UpdateAlbumRequest) *wrapper.ErrorCode {
	return a.updateAlbum.Handle(request)
}
