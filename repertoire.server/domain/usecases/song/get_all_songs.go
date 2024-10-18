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

func (g GetAllSongs) Handle(request requests.GetSongsRequest) (songs []models.Song, e *wrapper.ErrorCode) {
	err := g.repository.GetAllByUser(&songs, request.UserID, request.CurrentPage, request.PageSize)
	if err != nil {
		return songs, wrapper.InternalServerError(err)
	}
	return songs, nil
}
