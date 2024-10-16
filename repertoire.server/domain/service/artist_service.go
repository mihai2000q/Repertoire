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

func (s *artistService) Get(id uuid.UUID) (models.Artist, *utils.ErrorCode) {
	return s.getArtist.Handle(id)
}

func (s *artistService) GetAll(request requests.GetArtistsRequest) (artists []models.Artist, e *utils.ErrorCode) {
	return s.getAllArtists.Handle(request)
}

func (s *artistService) Create(request requests.CreateArtistRequest, token string) *utils.ErrorCode {
	return s.createArtist.Handle(request, token)
}

func (s *artistService) Update(request requests.UpdateArtistRequest) *utils.ErrorCode {
	return s.updateArtist.Handle(request)
}

func (s *artistService) Delete(id uuid.UUID) *utils.ErrorCode {
	return s.deleteArtist.Handle(id)
}
