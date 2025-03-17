package section

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateSongSectionsPartialOccurrences struct {
	songRepository repository.SongRepository
}

func NewUpdateSongSectionsPartialOccurrences(
	repository repository.SongRepository,
) UpdateSongSectionsPartialOccurrences {
	return UpdateSongSectionsPartialOccurrences{
		songRepository: repository,
	}
}

func (c UpdateSongSectionsPartialOccurrences) Handle(
	request requests.UpdateSongSectionsPartialOccurrencesRequest,
) *wrapper.ErrorCode {
	var song model.Song
	err := c.songRepository.GetWithSections(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	sectionsPartialOccurrencesMap := make(map[uuid.UUID]uint)
	for _, s := range request.Sections {
		sectionsPartialOccurrencesMap[s.ID] = s.PartialOccurrences
	}

	for i := range song.Sections {
		occurrences, ok := sectionsPartialOccurrencesMap[song.Sections[i].ID]
		if ok {
			song.Sections[i].PartialOccurrences = occurrences
		}
	}

	err = c.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
