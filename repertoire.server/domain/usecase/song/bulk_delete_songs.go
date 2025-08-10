package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"slices"
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
	var albumsToReorder []model.Album
	for _, song := range songs {
		if song.AlbumID != nil && !slices.ContainsFunc(albumsToReorder, func(album model.Album) bool {
			return album.ID == *song.AlbumID
		}) {
			albumsToReorder = append(albumsToReorder, *song.Album)
		}
	}

	var albumSongsToUpdate []model.Song
	for _, album := range albumsToReorder {
		songsFound := uint(0)
		for _, albumSong := range album.Songs {
			if slices.ContainsFunc(songs, func(song model.Song) bool {
				return song.ID == albumSong.ID
			}) {
				songsFound++
				continue
			}
			if songsFound != 0 {
				trackNo := *albumSong.AlbumTrackNo - songsFound
				albumSong.AlbumTrackNo = &trackNo
				albumSongsToUpdate = append(albumSongsToUpdate, albumSong)
			}
		}
	}

	if len(albumSongsToUpdate) != 0 {
		err := b.repository.UpdateAll(&albumSongsToUpdate)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	return nil
}

func (b BulkDeleteSongs) reorderSongsInPlaylists(songs []model.Song) *wrapper.ErrorCode {
	var playlistsToReorder []model.Playlist
	for _, song := range songs {
		for _, playlist := range song.Playlists {
			if !slices.ContainsFunc(playlistsToReorder, func(p model.Playlist) bool {
				return p.ID == playlist.ID
			}) {
				playlistsToReorder = append(playlistsToReorder, playlist)
			}
		}
	}

	var playlistSongsToUpdate []model.PlaylistSong
	for _, playlist := range playlistsToReorder {
		songsFound := uint(0)
		for _, playlistSong := range playlist.PlaylistSongs {
			if slices.ContainsFunc(songs, func(song model.Song) bool {
				return song.ID == playlistSong.SongID
			}) {
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
		err := b.playlistRepository.UpdateAllPlaylistSongs(&playlistSongsToUpdate)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}

	return nil
}
