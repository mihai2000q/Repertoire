package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateAllSongParts struct {
	songRepository repository.SongRepository
}

func NewUpdateAllSongParts(songRepository repository.SongRepository) UpdateAllSongParts {
	return UpdateAllSongParts{songRepository: songRepository}
}

func (u UpdateAllSongParts) Handle(request requests.UpdateAllSongPartsRequest) *wrapper.ErrorCode {
	var song model.Song
	err := u.songRepository.GetWithParts(&song, request.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	for i := range song.Parts {
		if request.InstrumentID != nil {
			song.Parts[i].InstrumentID = request.InstrumentID
		}
		if request.BandMemberID != nil {
			song.Parts[i].BandMemberID = request.BandMemberID
		}
	}

	err = u.songRepository.UpdateWithAssociations(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}
