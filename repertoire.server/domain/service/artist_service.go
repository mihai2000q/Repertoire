package service

import (
	"repertoire/api/request"
	"repertoire/domain/usecase/artist"
	"repertoire/model"
	"repertoire/utils/wrapper"

	"github.com/google/uuid"
)

type ArtistService interface {
	Get(id uuid.UUID) (model.Artist, *wrapper.ErrorCode)
	GetAll(request request.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode)
	Create(request request.CreateArtistRequest, token string) *wrapper.ErrorCode
	Update(request request.UpdateArtistRequest) *wrapper.ErrorCode
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

func (a *artistService) GetAll(request request.GetArtistsRequest, token string) (wrapper.WithTotalCount[model.Artist], *wrapper.ErrorCode) {
	return a.getAllArtists.Handle(request, token)
}

func (a *artistService) Create(request request.CreateArtistRequest, token string) *wrapper.ErrorCode {
	return a.createArtist.Handle(request, token)
}

func (a *artistService) Update(request request.UpdateArtistRequest) *wrapper.ErrorCode {
	return a.updateArtist.Handle(request)
}

func (a *artistService) Delete(id uuid.UUID) *wrapper.ErrorCode {
	return a.deleteArtist.Handle(id)
}
