package service

import (
	"mime/multipart"
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type ArtistService interface {
	Create(request requests.CreateArtistRequest, token string) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
	DeleteImage(id uuid.UUID) *wrapper.ErrorCode
	GetAll(request requests.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode)
	Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode)
	SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode
	Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode
}

type artistService struct {
	createArtist          artist.CreateArtist
	deleteArtist          artist.DeleteArtist
	deleteImageFromArtist artist.DeleteImageFromArtist
	getAllArtists         artist.GetAllArtists
	getArtist             artist.GetArtist
	saveImageToArtist     artist.SaveImageToArtist
	updateArtist          artist.UpdateArtist
}

func NewArtistService(
	createArtist artist.CreateArtist,
	deleteArtist artist.DeleteArtist,
	deleteImageFromArtist artist.DeleteImageFromArtist,
	getAllArtists artist.GetAllArtists,
	getArtist artist.GetArtist,
	saveImageToArtist artist.SaveImageToArtist,
	updateArtist artist.UpdateArtist,
) ArtistService {
	return &artistService{
		createArtist:          createArtist,
		deleteArtist:          deleteArtist,
		deleteImageFromArtist: deleteImageFromArtist,
		getAllArtists:         getAllArtists,
		getArtist:             getArtist,
		saveImageToArtist:     saveImageToArtist,
		updateArtist:          updateArtist,
	}
}

func (a *artistService) Create(request requests.CreateArtistRequest, token string) *wrapper.ErrorCode {
	return a.createArtist.Handle(request, token)
}

func (a *artistService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteArtist.Handle(id)
}

func (a *artistService) DeleteImage(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteImageFromArtist.Handle(id)
}

func (a *artistService) GetAll(request requests.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode) {
	return a.getAllArtists.Handle(request, token)
}

func (a *artistService) Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode) {
	return a.getArtist.Handle(id)
}

func (a *artistService) SaveImage(file *multipart.FileHeader, id uuid.UUID) *wrapper.ErrorCode {
	return a.saveImageToArtist.Handle(file, id)
}

func (a *artistService) Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode {
	return a.updateArtist.Handle(request)
}
