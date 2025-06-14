package album

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddSongsToAlbum struct {
	repository              repository.AlbumRepository
	songRepository          repository.SongRepository
	messagePublisherService service.MessagePublisherService
}

func NewAddSongsToAlbum(
	repository repository.AlbumRepository,
	songRepository repository.SongRepository,
	messagePublisherService service.MessagePublisherService,
) AddSongsToAlbum {
	return AddSongsToAlbum{
		repository:              repository,
		songRepository:          songRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (a AddSongsToAlbum) Handle(request requests.AddSongsToAlbumRequest) *wrapper.ErrorCode {
	var album model.Album
	err := a.repository.GetWithSongs(&album, request.ID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(album).IsZero() {
		return wrapper.NotFoundError(errors.New("album not found"))
	}

	var songs []model.Song
	err = a.songRepository.GetAllByIDs(&songs, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	songsLength := len(album.Songs) + 1
	for i, song := range songs {
		// if their artists don't match, or the song has an artist but the album doesn't, it results in failure
		// on the other hand, if the album has an artist and the song doesn't, it will inherit it (pass)
		if song.ArtistID != nil && (album.ArtistID == nil || *album.ArtistID != *song.ArtistID) {
			return wrapper.ConflictError(errors.New("song " + song.ID.String() + song.Title + " and album do not share the same artist"))
		}
		if song.AlbumID != nil {
			return wrapper.ConflictError(errors.New("song " + song.ID.String() + " already has an album"))
		}

		songs[i].AlbumID = &request.ID
		trackNo := uint(songsLength + i)
		songs[i].AlbumTrackNo = &trackNo
		songs[i].ArtistID = album.ArtistID // songs inherit album artist
		if songs[i].ReleaseDate == nil {
			songs[i].ReleaseDate = album.ReleaseDate // also inherit the release date if there is none
		}
	}

	err = a.songRepository.UpdateAll(&songs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = a.messagePublisherService.Publish(topics.SongsUpdatedTopic, request.SongIDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
