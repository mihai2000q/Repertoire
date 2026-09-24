package section

import (
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type GetAllSongSections struct {
	songSectionRepository repository.SongSectionRepository
}

func NewGetAllSongSections(songSectionRepository repository.SongSectionRepository) GetAllSongSections {
	return GetAllSongSections{songSectionRepository: songSectionRepository}
}

func (g GetAllSongSections) Handle(request requests.GetSongSectionsRequest) ([]model.SongSection, *httperror.ErrorCode) {
	var sections []model.SongSection
	err := g.songSectionRepository.GetAllBySong(&sections, request.SongID)
	if err != nil {
		return sections, httperror.DatabaseError(err)
	}
	return sections, nil
}
