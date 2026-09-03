package song

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"
	"sync"

	"github.com/google/uuid"
)

type DeleteSong struct {
	repository              repository.SongRepository
	playlistRepository      repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewDeleteSong(
	repository repository.SongRepository,
	playlistRepository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) DeleteSong {
	return DeleteSong{
		repository:              repository,
		playlistRepository:      playlistRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (d DeleteSong) Handle(id uuid.UUID) *httperror.ErrorCode {
	var song model.Song
	if err := d.repository.GetWithPlaylistsAndSongs(&song, id); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	var wg sync.WaitGroup
	errChan := make(chan *httperror.ErrorCode, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		errChan <- d.reorderAlbum(song)
	}()
	go func() {
		defer wg.Done()
		errChan <- d.reorderSongsInPlaylists(song)
	}()

	wg.Wait()
	close(errChan)
	for errorCode := range errChan {
		if errorCode != nil {
			return errorCode
		}
	}

	if err := d.repository.Delete([]uuid.UUID{id}); err != nil {
		return httperror.DatabaseError(err)
	}

	if err := d.messagePublisherService.Publish(topics.SongsDeletedTopic, []model.Song{song}); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}

func (d DeleteSong) reorderAlbum(song model.Song) *httperror.ErrorCode {
	if song.AlbumID == nil {
		return nil
	}

	var albumSongs []model.Song
	err := d.repository.GetAllByAlbumAndTrackNo(&albumSongs, *song.AlbumID, *song.AlbumTrackNo)
	if err != nil {
		return httperror.DatabaseError(err)
	}

	for i := range albumSongs {
		trackNo := *albumSongs[i].AlbumTrackNo - 1
		albumSongs[i].AlbumTrackNo = &trackNo
	}

	if len(albumSongs) != 0 {
		if err := d.repository.UpdateAll(&albumSongs); err != nil {
			return httperror.DatabaseError(err)
		}
	}

	return nil
}

func (d DeleteSong) reorderSongsInPlaylists(song model.Song) *httperror.ErrorCode {
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
		if err := d.playlistRepository.UpdateAllPlaylistSongs(&playlistSongsToUpdate); err != nil {
			return httperror.DatabaseError(err)
		}
	}

	return nil
}
