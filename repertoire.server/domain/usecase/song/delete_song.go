package song

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type DeleteSong struct {
	repository repository.SongRepository
}

func NewDeleteSong(repository repository.SongRepository) DeleteSong {
	return DeleteSong{repository: repository}
}

func (d DeleteSong) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.repository.Get(&song, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	if song.AlbumID != nil {
		errCode := d.reorderAlbum(song)
		if errCode != nil {
			return errCode
		}
	}

	err = d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	return nil
}

func (d DeleteSong) reorderAlbum(song model.Song) *wrapper.ErrorCode {
	var albumSongs []model.Song
	err := d.repository.GetAllByAlbumAndTrackNo(&albumSongs, *song.AlbumID, *song.AlbumTrackNo)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for i := range albumSongs {
		trackNo := *albumSongs[i].AlbumTrackNo - 1
		albumSongs[i].AlbumTrackNo = &trackNo
	}

	err = d.repository.UpdateAll(&albumSongs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
