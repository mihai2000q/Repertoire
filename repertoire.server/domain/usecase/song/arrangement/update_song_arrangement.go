package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UpdateSongArrangement struct {
	songArrangementRepository repository.SongArrangementRepository
}

func NewUpdateSongArrangement(songArrangementRepository repository.SongArrangementRepository) UpdateSongArrangement {
	return UpdateSongArrangement{songArrangementRepository: songArrangementRepository}
}

func (u UpdateSongArrangement) Handle(request requests.UpdateSongArrangementRequest) *wrapper.ErrorCode {
	var arrangement model.SongArrangement
	err := u.songArrangementRepository.GetWithAssociations(&arrangement, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(arrangement).IsZero() {
		return wrapper.NotFoundError(errors.New("song arrangement not found"))
	}

	arrangement.Name = request.Name

	// in case the sections in the request and from repository are not in the same order
	sectionsOccurrencesMap := make(map[uuid.UUID]uint)
	for _, s := range request.Occurrences {
		sectionsOccurrencesMap[s.SectionID] = s.Occurrences
	}

	// propagate the occurrences on the arrangement
	for i := range arrangement.Occurrences {
		occurrences, ok := sectionsOccurrencesMap[arrangement.Occurrences[i].SectionID]
		if ok {
			arrangement.Occurrences[i].Occurrences = occurrences
		}
	}

	err = u.songArrangementRepository.UpdateWithAssociations(&arrangement)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
