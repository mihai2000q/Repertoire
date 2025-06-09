package song

import (
	"errors"
	"github.com/google/uuid"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type UpdateSong struct {
	repository              repository.SongRepository
	albumRepository         repository.AlbumRepository
	messagePublisherService service.MessagePublisherService
}

func NewUpdateSong(
	repository repository.SongRepository,
	albumRepository repository.AlbumRepository,
	messagePublisherService service.MessagePublisherService,
) UpdateSong {
	return UpdateSong{
		repository:              repository,
		albumRepository:         albumRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (u UpdateSong) Handle(request requests.UpdateSongRequest) *wrapper.ErrorCode {
	var song model.Song
	err := u.repository.Get(&song, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	artistHasChanged := song.ArtistID != nil && request.ArtistID == nil ||
		song.ArtistID == nil && request.ArtistID != nil ||
		song.ArtistID != nil && request.ArtistID != nil && *song.ArtistID != *request.ArtistID
	albumHasChanged := song.AlbumID != nil && request.AlbumID == nil ||
		song.AlbumID == nil && request.AlbumID != nil ||
		song.AlbumID != nil && request.AlbumID != nil && *song.AlbumID != *request.AlbumID

	if (albumHasChanged || artistHasChanged) && request.AlbumID != nil {
		var album model.Album
		err = u.albumRepository.Get(&album, *request.AlbumID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		if reflect.ValueOf(album).IsZero() {
			return wrapper.NotFoundError(errors.New("album not found"))
		}
		if request.ArtistID == nil && album.ArtistID != nil ||
			request.ArtistID != nil && album.ArtistID == nil ||
			request.ArtistID != nil && album.ArtistID != nil && *request.ArtistID != *album.ArtistID {
			return wrapper.BadRequestError(errors.New("album's artist does not match the request's artist"))
		}
	}

	if albumHasChanged {
		errCode := u.reorderAlbumSongs(request, &song)
		if errCode != nil {
			return errCode
		}
	}

	song.Title = request.Title
	song.Description = request.Description
	song.IsRecorded = request.IsRecorded
	song.Bpm = request.Bpm
	song.SongsterrLink = request.SongsterrLink
	song.YoutubeLink = request.YoutubeLink
	song.ReleaseDate = request.ReleaseDate
	song.Difficulty = request.Difficulty
	song.GuitarTuningID = request.GuitarTuningID
	song.ArtistID = request.ArtistID
	song.AlbumID = request.AlbumID

	err = u.repository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = u.messagePublisherService.Publish(topics.SongsUpdatedTopic, []uuid.UUID{song.ID})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (u UpdateSong) reorderAlbumSongs(request requests.UpdateSongRequest, song *model.Song) *wrapper.ErrorCode {
	// reorder old album, if any
	if song.AlbumID != nil {
		var songs []model.Song
		err := u.repository.GetAllByAlbumAndTrackNo(&songs, *song.AlbumID, *song.AlbumTrackNo)
		if err != nil {
			return wrapper.InternalServerError(err)
		}

		for i := range songs {
			trackNo := *songs[i].AlbumTrackNo - 1
			songs[i].AlbumTrackNo = &trackNo
		}

		err = u.repository.UpdateAll(&songs)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	// add album track number on new song, if any new album
	if request.AlbumID == nil {
		return nil
	}

	var songsCount int64
	err := u.repository.CountByAlbum(&songsCount, *request.AlbumID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	trackNo := uint(songsCount) + 1
	song.AlbumTrackNo = &trackNo

	return nil
}
