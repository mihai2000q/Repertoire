package service

import (
	"repertoire/api/requests"
	"repertoire/domain/usecases/artist"
	"repertoire/models"
	"repertoire/utils"

	"github.com/google/uuid"
)

type ArtistService interface {
	Get(id uuid.UUID) (artist models.Artist, e *utils.ErrorCode)
	GetAll(request requests.GetArtistsRequest) (artists []models.Artist, e *utils.ErrorCode)
	Create(request requests.CreateArtistRequest, token string) *utils.ErrorCode
	Update(request requests.UpdateArtistRequest) *utils.ErrorCode
	Delete(id uuid.UUID) *utils.ErrorCode
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

func (a *artistService) Get(id uuid.UUID) (models.Artist, *utils.ErrorCode) {
	return a.getArtist.Handle(id)
}

func (a *artistService) GetAll(request requests.GetArtistsRequest) ([]models.Artist, *utils.ErrorCode) {
	return a.getAllArtists.Handle(request)
}

func (a *artistService) Create(request requests.CreateArtistRequest, token string) *utils.ErrorCode {
	return a.createArtist.Handle(request, token)
}

func (a *artistService) Update(request requests.UpdateArtistRequest) *utils.ErrorCode {
	return a.updateArtist.Handle(request)
}

func (a *artistService) Delete(id uuid.UUID) *utils.ErrorCode {
	return a.deleteArtist.Handle(id)
}
