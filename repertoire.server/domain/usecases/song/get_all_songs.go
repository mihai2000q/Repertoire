package song

import (
	"repertoire/api/requests"
	"repertoire/data/repository"
	"repertoire/models"
	"repertoire/utils"
)

type GetAllSongs struct {
	repository repository.SongRepository
}

func NewGetAllSongs(repository repository.SongRepository) *GetAllSongs {
	return &GetAllSongs{
		repository: repository,
	}
}

func (g GetAllSongs) Handle(request requests.GetSongsRequest) (songs []models.Song, e *utils.ErrorCode) {
	err := g.repository.GetAllByUser(&songs, request.UserID)
	if err != nil {
		return songs, utils.InternalServerError(err)
	}
	return songs, nil
}
