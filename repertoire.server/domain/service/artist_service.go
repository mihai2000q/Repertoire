package service

import (
	"repertoire/api/requests"
	"repertoire/domain/usecases/artist"
	"repertoire/models"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
)

type ArtistService interface {
	Get(id uuid.UUID) (models.Artist, *wrapper.ErrorCode)
	GetAll(request requests.GetArtistsRequest) (wrapper.WithTotalCount[models.Artist], *wrapper.ErrorCode)
	Create(request requests.CreateArtistRequest, token string) *wrapper.ErrorCode
	Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode
	Delete(id uuid.UUID) *wrapper.ErrorCode
}

type artistService struct {
	getArtist     artist.GetArtist
	getAllArtists artist.GetAllArtists
	createArtist  artist.CreateArtist
	updateArtist  artist.UpdateArtist
	deleteArtist  artist.DeleteArtist
}

func NewArtistService(
	getArtist artist.GetArtist,
	getAllArtists artist.GetAllArtists,
	createArtist artist.CreateArtist,
	updateArtist artist.UpdateArtist,
	deleteArtist artist.DeleteArtist,
) ArtistService {
	return &artistService{
		getArtist:     getArtist,
		getAllArtists: getAllArtists,
		createArtist:  createArtist,
		updateArtist:  updateArtist,
		deleteArtist:  deleteArtist,
	}
}

func (a *artistService) Get(id uuid.UUID) (models.Artist, *wrapper.ErrorCode) {
	return a.getArtist.Handle(id)
}

func (a *artistService) GetAll(request requests.GetArtistsRequest) (wrapper.WithTotalCount[models.Artist], *wrapper.ErrorCode) {
	return a.getAllArtists.Handle(request)
}

func (a *artistService) Create(request requests.CreateArtistRequest, token string) *wrapper.ErrorCode {
	return a.createArtist.Handle(request, token)
}

func (a *artistService) Update(request requests.UpdateArtistRequest) *wrapper.ErrorCode {
	return a.updateArtist.Handle(request)
}

func (a *artistService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteArtist.Handle(id)
}
