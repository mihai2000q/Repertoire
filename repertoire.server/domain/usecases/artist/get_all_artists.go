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

func (g GetAllArtists) Handle(request requests.GetArtistsRequest) (result wrapper.WithTotalCount[models.Artist], e *wrapper.ErrorCode) {
	err := g.repository.GetAllByUser(&result.Data, request.UserID, request.CurrentPage, request.PageSize)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}
	err = g.repository.GetAllByUserCount(&result.TotalCount, request.UserID)
	if err != nil {
		return result, wrapper.InternalServerError(err)
	}
	return result, nil
}
