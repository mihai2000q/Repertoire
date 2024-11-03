package album

import (
	"errors"
	"github.com/google/uuid"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type MoveSongFromAlbum struct {
	repository repository.AlbumRepository
}

func NewMoveSongFromAlbum(repository repository.AlbumRepository) MoveSongFromAlbum {
	return MoveSongFromAlbum{repository: repository}
}

func (m MoveSongFromAlbum) Handle(request requests.MoveSongFromAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := m.repository.GetWithSongs(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if album.ID == uuid.Nil {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	index, overIndex, err := m.getIndexes(album.Songs, request.SongID, request.OverSongID)
	if err != nil {
		return wrapper.NotFoundError(err)
	}
	album.Songs = m.move(album.Songs, index, overIndex)

	err = m.repository.UpdateWithAssociations(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (MoveSongFromAlbum) getIndexes(songs []model.Song, id uuid.UUID, overID uuid.UUID) (int, int, error) {
	var index *int
	var overIndex *int
	for i := 0; i < len(songs); i++ {
		if songs[i].ID == id {
			index = &i
		} else if songs[i].ID == overID {
			overIndex = &i
		}
	}

	if index == nil {
		return -1, -1, errors.New("song not found")
	}
	if overIndex == nil {
		return -1, -1, errors.New("over song not found")
	}

	return *index, *overIndex, nil
}

func (MoveSongFromAlbum) move(songs []model.Song, index int, overIndex int) []model.Song {
	if index < overIndex {
		for i := index + 1; i <= overIndex; i++ {
			trackNo := *songs[i].AlbumTrackNo - 1
			songs[i].AlbumTrackNo = &trackNo
		}
	} else {
		for i := overIndex; i <= index; i++ {
			trackNo := *songs[i].AlbumTrackNo + 1
			songs[i].AlbumTrackNo = &trackNo
		}
	}

	trackNo := uint(overIndex) + 1
	songs[index].AlbumTrackNo = &trackNo

	return songs
}
