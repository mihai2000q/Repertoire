package artist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type GetAllArtists struct {
	repository repository.ArtistRepository
}

func NewGetAllArtists(repository repository.ArtistRepository) GetAllArtists {
	return GetAllArtists{
		repository: repository,
	}
}

func (g GetAllArtists) Handle(request requests.GetArtistsRequest) (artists []models.Artist, e *utils.ErrorCode) {
	err := g.repository.GetAllByUser(&artists, request.UserID)
	if err != nil {
		return artists, utils.InternalServerError(err)
	}
	return artists, nil
}
