package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"sync"
)

type BulkDeleteSongs struct {
	repository              repository.SongRepository
	playlistRepository      repository.PlaylistRepository
	messagePublisherService service.MessagePublisherService
}

func NewBulkDeleteSongs(
	repository repository.SongRepository,
	playlistRepository repository.PlaylistRepository,
	messagePublisherService service.MessagePublisherService,
) BulkDeleteSongs {
	return BulkDeleteSongs{
		repository:              repository,
		playlistRepository:      playlistRepository,
		messagePublisherService: messagePublisherService,
	}
}

func (b BulkDeleteSongs) Handle(request requests.BulkDeleteSongsRequest) *wrapper.ErrorCode {
	var songs []model.Song
	err := b.repository.GetAllByIDsWithAlbumsAndPlaylists(&songs, request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(songs) == 0 {
		return wrapper.NotFoundError(errors.New("songs not found"))
	}

	var wg sync.WaitGroup
	errChan := make(chan *wrapper.ErrorCode, 2)
	wg.Add(2)
	go func() {
		defer wg.Done()
		errChan <- b.reorderAlbums(songs)
	}()
	go func() {
		defer wg.Done()
		errChan <- b.reorderSongsInPlaylists(songs)
	}()

	wg.Wait()
	close(errChan)
	for errorCode := range errChan {
		if errorCode != nil {
			return errorCode
		}
	}

	err = b.repository.Delete(request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = b.messagePublisherService.Publish(topics.SongsDeletedTopic, songs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (b BulkDeleteSongs) reorderAlbums(songs []model.Song) *wrapper.ErrorCode {
	var albumSongsToUpdate []model.Song
	for _, song := range songs {
		if song.AlbumID == nil {
			continue
		}

		// PROBLEM:
		// I delete a song from songs[0].Album.Songs
		// but songs[1].Album.Songs has the same Album as songs[0],
		// and on songs[1].Album.Songs, songs[0] still exists.
		songFound := false
		for _, albumSong := range song.Album.Songs {
			if albumSong.ID == song.ID {
				songFound = true
				continue
			}
			if songFound {
				trackNo := *albumSong.AlbumTrackNo - 1
				albumSong.AlbumTrackNo = &trackNo
				albumSongsToUpdate = append(albumSongsToUpdate, albumSong)
			}
		}
	}

	if len(albumSongsToUpdate) == 0 {
		return nil
	}
	err := b.repository.UpdateAll(&albumSongsToUpdate)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (b BulkDeleteSongs) reorderSongsInPlaylists(songs []model.Song) *wrapper.ErrorCode {
	var playlistSongsToUpdate []model.PlaylistSong
	for _, song := range songs {
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
	}

	if len(playlistSongsToUpdate) == 0 {
		return nil
	}
	err := b.playlistRepository.UpdateAllPlaylistSongs(&playlistSongsToUpdate)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
