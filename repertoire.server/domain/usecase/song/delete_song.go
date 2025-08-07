package song

import (
	"errors"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
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

func (d DeleteSong) Handle(id uuid.UUID) *wrapper.ErrorCode {
	var song model.Song
	err := d.repository.GetWithPlaylistsAndSongs(&song, id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	var wg sync.WaitGroup
	errChan := make(chan *wrapper.ErrorCode, 2)
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

	err = d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.SongsDeletedTopic, []model.Song{song})
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (d DeleteSong) reorderAlbum(song model.Song) *wrapper.ErrorCode {
	if song.AlbumID == nil {
		return nil
	}

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

func (d DeleteSong) reorderSongsInPlaylists(song model.Song) *wrapper.ErrorCode {
	for _, playlist := range song.Playlists {
		songsFound := uint(0)
		var playlistSongsToUpdate []model.PlaylistSong
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
		err := d.playlistRepository.UpdateAllPlaylistSongs(&playlistSongsToUpdate)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}
	return nil
}
