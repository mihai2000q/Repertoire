package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/album"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/pagination"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type AlbumService interface {
	AddSongs(requests.AddSongsToAlbumRequest) *httperror.ErrorCode
	AddPerfectRehearsals(request requests.AddPerfectRehearsalsToAlbumsRequest) *httperror.ErrorCode
	BulkDelete(request requests.BulkDeleteAlbumsRequest) *httperror.ErrorCode
	Create(request requests.CreateAlbumRequest, token string) (uuid.UUID, *httperror.ErrorCode)
	Delete(request requests.DeleteAlbumRequest) *httperror.ErrorCode
	DeleteImage(id uuid.UUID) *httperror.ErrorCode
	Get(request requests.GetAlbumRequest) (model.Album, *httperror.ErrorCode)
	GetFiltersMetadata(
		request requests.GetAlbumFiltersMetadataRequest,
		token string,
	) (model.AlbumFiltersMetadata, *httperror.ErrorCode)
	GetAll(request requests.GetAlbumsRequest, token string) (pagination.WithTotalCount[model.EnhancedAlbum], *httperror.ErrorCode)
	MoveSong(request requests.MoveSongFromAlbumRequest) *httperror.ErrorCode
	RemoveSongs(request requests.RemoveSongsFromAlbumRequest) *httperror.ErrorCode
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode
	Update(request requests.UpdateAlbumRequest) *httperror.ErrorCode
}

type albumService struct {
	addSongsToAlbum              album.AddSongsToAlbum
	addPerfectRehearsalsToAlbums album.AddPerfectRehearsalsToAlbums
	bulkDeleteAlbums             album.BulkDeleteAlbums
	createAlbum                  album.CreateAlbum
	deleteAlbum                  album.DeleteAlbum
	deleteImageFromAlbum         album.DeleteImageFromAlbum
	getAlbum                     album.GetAlbum
	getAlbumFiltersMetadata      album.GetAlbumFiltersMetadata
	getAllAlbums                 album.GetAllAlbums
	moveSongFromAlbum            album.MoveSongFromAlbum
	removeSongsFromAlbum         album.RemoveSongsFromAlbum
	saveImageToAlbum             album.SaveImageToAlbum
	updateAlbum                  album.UpdateAlbum
}

func NewAlbumService(
	addSongsToAlbum album.AddSongsToAlbum,
	addPerfectRehearsalsToAlbums album.AddPerfectRehearsalsToAlbums,
	bulkDeleteAlbums album.BulkDeleteAlbums,
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
		addSongsToAlbum:              addSongsToAlbum,
		addPerfectRehearsalsToAlbums: addPerfectRehearsalsToAlbums,
		bulkDeleteAlbums:             bulkDeleteAlbums,
		createAlbum:                  createAlbum,
		deleteAlbum:                  deleteAlbum,
		deleteImageFromAlbum:         deleteImageFromAlbum,
		getAlbum:                     getAlbum,
		getAlbumFiltersMetadata:      getAlbumFiltersMetadata,
		getAllAlbums:                 getAllAlbums,
		moveSongFromAlbum:            moveSongFromAlbum,
		removeSongsFromAlbum:         removeSongsFromAlbum,
		saveImageToAlbum:             saveImageToAlbum,
		updateAlbum:                  updateAlbum,
	}
}

func (a *albumService) AddSongs(request requests.AddSongsToAlbumRequest) *httperror.ErrorCode {
	return a.addSongsToAlbum.Handle(request)
}

func (a *albumService) AddPerfectRehearsals(request requests.AddPerfectRehearsalsToAlbumsRequest) *httperror.ErrorCode {
	return a.addPerfectRehearsalsToAlbums.Handle(request)
}

func (a *albumService) BulkDelete(request requests.BulkDeleteAlbumsRequest) *httperror.ErrorCode {
	return a.bulkDeleteAlbums.Handle(request)
}

func (a *albumService) Create(request requests.CreateAlbumRequest, token string) (uuid.UUID, *httperror.ErrorCode) {
	return a.createAlbum.Handle(request, token)
}

func (a *albumService) Delete(request requests.DeleteAlbumRequest) *httperror.ErrorCode {
	return a.deleteAlbum.Handle(request)
}

func (a *albumService) DeleteImage(id uuid.UUID) *httperror.ErrorCode {
	return a.deleteImageFromAlbum.Handle(id)
}

func (a *albumService) Get(request requests.GetAlbumRequest) (model.Album, *httperror.ErrorCode) {
	return a.getAlbum.Handle(request)
}

func (a *albumService) GetFiltersMetadata(
	request requests.GetAlbumFiltersMetadataRequest,
	token string,
) (model.AlbumFiltersMetadata, *httperror.ErrorCode) {
	return a.getAlbumFiltersMetadata.Handle(request, token)
}

func (a *albumService) GetAll(request requests.GetAlbumsRequest, token string) (pagination.WithTotalCount[model.EnhancedAlbum], *httperror.ErrorCode) {
	return a.getAllAlbums.Handle(request, token)
}

func (a *albumService) MoveSong(request requests.MoveSongFromAlbumRequest) *httperror.ErrorCode {
	return a.moveSongFromAlbum.Handle(request)
}

func (a *albumService) RemoveSongs(request requests.RemoveSongsFromAlbumRequest) *httperror.ErrorCode {
	return a.removeSongsFromAlbum.Handle(request)
}

func (a *albumService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *httperror.ErrorCode {
	return a.saveImageToAlbum.Handle(file, id)
}

func (a *albumService) Update(request requests.UpdateAlbumRequest) *httperror.ErrorCode {
	return a.updateAlbum.Handle(request)
}
