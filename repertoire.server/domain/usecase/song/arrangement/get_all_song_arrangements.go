package arrangement

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type GetAllSongArrangements struct {
	songArrangementRepository repository.SongArrangementRepository
}

func NewGetAllSongArrangements(songArrangementRepository repository.SongArrangementRepository) GetAllSongArrangements {
	return GetAllSongArrangements{songArrangementRepository: songArrangementRepository}
}

func (g GetAllSongArrangements) Handle(request requests.GetSongArrangementsRequest) ([]model.SongArrangement, *wrapper.ErrorCode) {
	var arrangements []model.SongArrangement
	err := g.songArrangementRepository.GetAllBySong(&arrangements, request.SongID)
	if err != nil {
		return arrangements, wrapper.InternalServerError(err)
	}
	return arrangements, nil
}
