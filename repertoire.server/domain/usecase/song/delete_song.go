package song

import (
	"errors"
	"reflect"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type DeleteSong struct {
	songRepository          repository.SongRepository
	transactionManager      transaction.Manager
	messagePublisherService service.MessagePublisherService

	txSongRepo     repository.SongRepository
	txPlaylistRepo repository.PlaylistRepository
}

func NewDeleteSong(
	songRepository repository.SongRepository,
	transactionManager transaction.Manager,
	messagePublisherService service.MessagePublisherService,
) DeleteSong {
	return DeleteSong{
		songRepository:          songRepository,
		transactionManager:      transactionManager,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteSong) Handle(id uuid.UUID) *httperror.ErrorCode {
	var song model.Song
	if err := d.songRepository.GetWithPlaylistsAndSongs(&song, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	err := d.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		d.txSongRepo = factory.NewSongRepository()
		d.txPlaylistRepo = factory.NewPlaylistRepository()

		if err := d.reorderAlbum(song); err != nil {
			return err
		}
		if err := d.reorderSongsInPlaylists(song); err != nil {
			return err
		}
		if err := d.txSongRepo.Delete([]uuid.UUID{id}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.SongsDeletedTopic, []model.Song{song}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}

func (d DeleteSong) reorderAlbum(song model.Song) error {
	if song.AlbumID == nil {
		return nil
	}

	var albumSongs []model.Song
	err := d.txSongRepo.GetAllByAlbumAndTrackNo(&albumSongs, *song.AlbumID, *song.AlbumTrackNo)
	if err != nil {
		return err
	}
	if len(albumSongs) == 0 {
		return nil
	}

	for i := range albumSongs {
		trackNo := *albumSongs[i].AlbumTrackNo - 1
		albumSongs[i].AlbumTrackNo = &trackNo
	}

	if err := d.txSongRepo.UpdateAll(&albumSongs); err != nil {
		return err
	}

	return nil
}

func (d DeleteSong) reorderSongsInPlaylists(song model.Song) error {
	var playlistSongsToUpdate []model.PlaylistSong
	for _, playlist := range song.Playlists {
		songsFound := uint(0)
		for _, playlistSong := range playlist.PlaylistSongs {
			if playlistSong.SongID == song.ID {
				songsFound++
				continue
			}

			if songsFound != 0 {
				playlistSong.SongTrackNo = playlistSong.SongTrackNo - songsFound
				playlistSongsToUpdate = append(playlistSongsToUpdate, playlistSong)
			}
		}
	}

	if len(playlistSongsToUpdate) != 0 {
		if err := d.txPlaylistRepo.UpdateAllPlaylistSongs(&playlistSongsToUpdate); err != nil {
			return err
		}
	}

	return nil
}
