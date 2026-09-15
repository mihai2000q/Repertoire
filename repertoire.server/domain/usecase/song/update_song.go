package song

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UpdateSong struct {
	songRepository          repository.SongRepository
	albumRepository         repository.AlbumRepository
	transactionManager      transaction.Manager
	messagePublisherService service.MessagePublisherService

	txSongRepo repository.SongRepository
}

func NewUpdateSong(
	songRepository repository.SongRepository,
	albumRepository repository.AlbumRepository,
	transactionManager transaction.Manager,
	messagePublisherService service.MessagePublisherService,
) UpdateSong {
	return UpdateSong{
		songRepository:          songRepository,
		albumRepository:         albumRepository,
		transactionManager:      transactionManager,
		messagePublisherService: messagePublisherService,
	}
}

func (u UpdateSong) Handle(request requests.UpdateSongRequest) *httperror.ErrorCode {
	var song model.Song
	if err := u.songRepository.Get(&song, request.ID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	albumHasChanged, errCode := u.ensureRequestArtistBelongsToAlbum(song, request)
	if errCode != nil {
		return errCode
	}

	err := u.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		u.txSongRepo = factory.NewSongRepository()

		if albumHasChanged {
			if err := u.reorderAlbumSongs(request, &song); err != nil {
				return err
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

		if err := u.txSongRepo.Update(&song); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	if err := u.messagePublisherService.Publish(topics.SongsUpdatedTopic, []uuid.UUID{song.ID}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}

func (u UpdateSong) ensureRequestArtistBelongsToAlbum(song model.Song, request requests.UpdateSongRequest) (bool, *httperror.ErrorCode) {
	albumHasChanged := song.AlbumID != nil && request.AlbumID == nil ||
		song.AlbumID == nil && request.AlbumID != nil ||
		song.AlbumID != nil && request.AlbumID != nil && *song.AlbumID != *request.AlbumID

	artistHasChanged := song.ArtistID != nil && request.ArtistID == nil ||
		song.ArtistID == nil && request.ArtistID != nil ||
		song.ArtistID != nil && request.ArtistID != nil && *song.ArtistID != *request.ArtistID

	if (albumHasChanged || artistHasChanged) && request.AlbumID != nil {
		var album model.Album
		if err := u.albumRepository.Get(&album, *request.AlbumID); err != nil {
			return false, httperror.DatabaseError(err)
		}
		if reflect.ValueOf(album).IsZero() {
			return false, httperror.NotFoundError(errors.New("album not found"))
		}
		if request.ArtistID == nil && album.ArtistID != nil ||
			request.ArtistID != nil && album.ArtistID == nil ||
			request.ArtistID != nil && album.ArtistID != nil && *request.ArtistID != *album.ArtistID {
			return false, httperror.ConflictError(errors.New("album's artist does not match the request's artist"))
		}
	}

	return albumHasChanged, nil
}

func (u UpdateSong) reorderAlbumSongs(request requests.UpdateSongRequest, song *model.Song) error {
	// reorder old album, if any
	if song.AlbumID != nil {
		var songs []model.Song
		err := u.txSongRepo.GetAllByAlbumAndTrackNo(&songs, *song.AlbumID, *song.AlbumTrackNo)
		if err != nil {
			return err
		}

		for i := range songs {
			trackNo := *songs[i].AlbumTrackNo - 1
			songs[i].AlbumTrackNo = &trackNo
		}

		if err := u.txSongRepo.UpdateAll(&songs); err != nil {
			return err
		}
	}

	// add album track number on new song, if any new album
	if request.AlbumID == nil {
		return nil
	}

	var songsCount int64
	if err := u.txSongRepo.CountByAlbum(&songsCount, *request.AlbumID); err != nil {
		return err
	}

	trackNo := uint(songsCount) + 1
	song.AlbumTrackNo = &trackNo

	return nil
}
