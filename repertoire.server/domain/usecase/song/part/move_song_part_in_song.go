package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/reorder"
	"repertoire/server/model"
)

type MoveSongPartInSong struct {
	songRepository repository.SongRepository
}

func NewMoveSongPartInSong(songRepository repository.SongRepository) MoveSongPartInSong {
	return MoveSongPartInSong{
		songRepository: songRepository,
	}
}

func (c MoveSongPartInSong) Handle(request requests.MoveSongPartInSongRequest) *httperror.ErrorCode {
	var song model.Song
	if err := c.songRepository.GetWithParts(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	errCode := reorder.MoveEntity(
		song.Parts,
		request.ID,
		request.OverID,
		&reorder.Config{
			OrderField:            "SongOrder",
			EntityNotFoundMsg:     "part not found",
			OverEntityNotFoundMsg: "over part not found",
		},
	)
	if errCode != nil {
		return errCode
	}

	if err := c.songRepository.UpdateWithAssociations(&song); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
