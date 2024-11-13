package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
)

type RemoveSongsFromAlbum struct {
	repository repository.AlbumRepository
}

func NewRemoveSongsFromAlbum(repository repository.AlbumRepository) RemoveSongsFromAlbum {
	return RemoveSongsFromAlbum{repository: repository}
}

func (r RemoveSongsFromAlbum) Handle(request requests.RemoveSongsFromAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := r.repository.GetWithSongs(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	var songsToDelete []model.Song
	var songsToPreserve []model.Song

	albumTrackNo := uint(1)
	for _, song := range album.Songs {
		if slices.Contains(request.SongIDs, song.ID) {
			songsToDelete = append(songsToDelete, song)
		} else {
			// reorder preserved songs
			*song.AlbumTrackNo = albumTrackNo
			songsToPreserve = append(songsToPreserve, song)
			albumTrackNo++
		}
	}

	if len(songsToDelete) != len(request.SongIDs) {
		return wrapper.NotFoundError(errors.New("could not find all songs"))
	}
	err = r.repository.RemoveSongs(&album, &songsToDelete)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	album.Songs = songsToPreserve
	err = r.repository.UpdateWithAssociations(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
