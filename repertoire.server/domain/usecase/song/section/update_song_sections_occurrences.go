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

type UpdateSongSectionsOccurrences struct {
	songRepository repository.SongRepository
}

func NewUpdateSongSectionsOccurrences(
	repository repository.SongRepository,
) UpdateSongSectionsOccurrences {
	return UpdateSongSectionsOccurrences{
		songRepository: repository,
	}
}

func (c UpdateSongSectionsOccurrences) Handle(
	request requests.UpdateSongSectionsOccurrencesRequest,
) *wrapper.ErrorCode {
	var song model.Song
	err := c.songRepository.GetWithSections(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	sectionsOccurrencesMap := make(map[uuid.UUID]uint)
	for _, s := range request.Sections {
		sectionsOccurrencesMap[s.ID] = s.Occurrences
	}

	for i := range song.Sections {
		occurrences, ok := sectionsOccurrencesMap[song.Sections[i].ID]
		if ok {
			song.Sections[i].Occurrences = occurrences
		}
	}

	err = c.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
