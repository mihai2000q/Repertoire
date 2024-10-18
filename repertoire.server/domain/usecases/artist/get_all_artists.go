package artist

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type GetAllArtists struct {
	repository repository.ArtistRepository
}

func NewGetAllArtists(repository repository.ArtistRepository) GetAllArtists {
	return GetAllArtists{
		repository: repository,
	}
}

func (g GetAllArtists) Handle(request requests.GetArtistsRequest) (artists []models.Artist, e *wrapper.ErrorCode) {
	err := g.repository.GetAllByUser(&artists, request.UserID, request.CurrentPage, request.PageSize)
	if err != nil {
		return artists, wrapper.InternalServerError(err)
	}
	return artists, nil
}
