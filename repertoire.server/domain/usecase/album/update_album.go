package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateAlbum struct {
	repository     repository.AlbumRepository
	songRepository repository.SongRepository
}

func NewUpdateAlbum(
	repository repository.AlbumRepository,
	songRepository repository.SongRepository,
) UpdateAlbum {
	return UpdateAlbum{
		repository:     repository,
		songRepository: songRepository,
	}
}

func (u UpdateAlbum) Handle(request requests.UpdateAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := u.repository.Get(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	artistHasChanged := album.ArtistID != request.ArtistID

	album.Title = request.Title
	album.ReleaseDate = request.ReleaseDate
	album.ArtistID = request.ArtistID

	err = u.repository.Update(&album)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	if artistHasChanged {
		errCode := u.updateAlbumSongsArtist(request)
		if errCode != nil {
			return errCode
		}
	}

	return nil
}

func (u UpdateAlbum) updateAlbumSongsArtist(request requests.UpdateAlbumRequest) *wrapper.ErrorCode {
	var songs []model.Song
	err := u.songRepository.GetAllByAlbum(&songs, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for i := range songs {
		songs[i].ArtistID = request.ArtistID
	}
	err = u.songRepository.UpdateAll(&songs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
