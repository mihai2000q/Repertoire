package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateAllSongSections struct {
	songRepository repository.SongRepository
}

func NewUpdateAllSongSections(songRepository repository.SongRepository) UpdateAllSongSections {
	return UpdateAllSongSections{songRepository: songRepository}
}

func (u UpdateAllSongSections) Handle(request requests.UpdateAllSongSectionsRequest) *wrapper.ErrorCode {
	var song model.Song
	err := u.songRepository.GetWithSections(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	for i := range song.Sections {
		if request.InstrumentID != nil {
			song.Sections[i].InstrumentID = request.InstrumentID
		}
		if request.BandMemberID != nil {
			song.Sections[i].BandMemberID = request.BandMemberID
		}
	}

	err = u.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
