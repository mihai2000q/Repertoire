package arrangement

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BulkUpdateSongArrangements struct {
	songArrangementRepository repository.SongArrangementRepository
}

func NewBulkUpdateSongArrangements(songArrangementRepository repository.SongArrangementRepository) BulkUpdateSongArrangements {
	return BulkUpdateSongArrangements{songArrangementRepository: songArrangementRepository}
}

func (b BulkUpdateSongArrangements) Handle(request requests.BulkUpdateSongArrangementsRequest) *wrapper.ErrorCode {
	// the arrangements in the request are not in the same order as those returned
	requestsMap := make(map[uuid.UUID]requests.UpdateSongArrangementRequest)
	var ids []uuid.UUID
	for _, r := range request.Requests {
		ids = append(ids, r.ID)
		requestsMap[r.ID] = r
	}

	var arrangements []model.SongArrangement
	err := b.songArrangementRepository.GetAllBySongWithPartOccurrences(&arrangements, ids, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(arrangements) != len(ids) {
		return wrapper.NotFoundError(errors.New("song arrangements not found"))
	}

	for i, arrangement := range arrangements {
		arrangements[i].Name = requestsMap[arrangement.ID].Name

		// in case the sections in the request and from repository are not in the same order
		partsOccurrencesMap := make(map[uuid.UUID]uint)
		for _, s := range requestsMap[arrangement.ID].Occurrences {
			partsOccurrencesMap[s.PartID] = s.Occurrences
		}

		// propagate the occurrences on the arrangement
		for j := range arrangement.PartOccurrences {
			occurrences, ok := partsOccurrencesMap[arrangement.PartOccurrences[j].PartID]
			if ok {
				arrangements[i].PartOccurrences[j].Occurrences = occurrences
			}
		}
	}

	err = b.songArrangementRepository.UpdateAllWithAssociations(&arrangements)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
