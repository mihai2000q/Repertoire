package song

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils/wrapper"
)

type GetAllSongs struct {
	repository repository.SongRepository
}

func NewGetAllSongs(repository repository.SongRepository) GetAllSongs {
	return GetAllSongs{
		repository: repository,
	}
}

func (g GetAllSongs) Handle(request requests.GetSongsRequest) (result wrapper.WithTotalCount[models.Song], e *wrapper.ErrorCode) {
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
