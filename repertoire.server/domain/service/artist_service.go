package service

import (
	"repertoire/server/api/requests"
	"repertoire/server/domain/usecase/artist"
	"repertoire/server/model"
	"repertoire/server/utils/wrapper"

	"github.com/google/uuid"
)

type ArtistService interface {
	Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode)
	GetAll(request requests.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode)
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

func (a *artistService) Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode) {
	return a.getArtist.Handle(id)
}

func (a *artistService) GetAll(request requests.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode) {
	return a.getAllArtists.Handle(request, token)
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
