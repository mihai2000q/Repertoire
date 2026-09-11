package part

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetAllSongParts struct {
	songPartRepository repository.SongPartRepository
}

func NewGetAllSongParts(songPartRepository repository.SongPartRepository) GetAllSongParts {
	return GetAllSongParts{songPartRepository: songPartRepository}
}

func (g GetAllSongParts) Handle(request requests.GetSongPartsRequest) ([]model.SongPart, *httperror.ErrorCode) {
	var parts []model.SongPart
	err := g.songPartRepository.GetAllBySong(&parts, request.SongID)
	if err != nil {
		return parts, httperror.DatabaseError(err)
	}
	return parts, nil
}
