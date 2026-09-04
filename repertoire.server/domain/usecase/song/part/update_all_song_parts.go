package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type UpdateAllSongParts struct {
	songRepository repository.SongRepository
}

func NewUpdateAllSongParts(songRepository repository.SongRepository) UpdateAllSongParts {
	return UpdateAllSongParts{songRepository: songRepository}
}

func (u UpdateAllSongParts) Handle(request requests.UpdateAllSongPartsRequest) *httperror.ErrorCode {
	var song model.Song
	if err := u.songRepository.GetWithParts(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	for i := range song.Parts {
		if request.InstrumentID != nil {
			song.Parts[i].InstrumentID = request.InstrumentID
		}
		if request.BandMemberID != nil {
			song.Parts[i].BandMemberID = request.BandMemberID
		}
	}

	if err := u.songRepository.UpdateWithAssociations(&song); err != nil {
		return httperror.DatabaseError(err)
	}
	return nil
}
