package arrangement

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type UpdateDefaultSongArrangement struct {
	songRepository repository.SongRepository
}

func NewUpdateDefaultSongArrangement(songRepository repository.SongRepository) UpdateDefaultSongArrangement {
	return UpdateDefaultSongArrangement{songRepository: songRepository}
}

func (g UpdateDefaultSongArrangement) Handle(request requests.UpdateDefaultSongArrangementRequest) *httperror.ErrorCode {
	var song model.Song
	if err := g.songRepository.Get(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	song.DefaultArrangementID = request.ID

	if err := g.songRepository.Update(&song); err != nil {
		return httperror.DatabaseError(err)
	}
	return nil
}
