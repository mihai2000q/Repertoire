package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/internal/deduplicate"
	"repertoire/server/internal/httperror"
	"repertoire/server/internal/message/topics"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type BulkDeleteSongs struct {
	songRepository          repository.SongRepository
	transactionManager      transaction.Manager
	messagePublisherService service.MessagePublisherService

	txSongRepo     repository.SongRepository
	txPlaylistRepo repository.PlaylistRepository
}

func NewBulkDeleteSongs(
	songRepository repository.SongRepository,
	transactionManager transaction.Manager,
	messagePublisherService service.MessagePublisherService,
) BulkDeleteSongs {
	return BulkDeleteSongs{
		songRepository:          songRepository,
		transactionManager:      transactionManager,
		messagePublisherService: messagePublisherService,
	}
}

func (b BulkDeleteSongs) Handle(request requests.BulkDeleteSongsRequest) *httperror.ErrorCode {
	var songs []model.Song
	if err := b.songRepository.GetAllByIDsWithAlbumsAndPlaylists(&songs, request.IDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(songs) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("songs not found"))
	}

	err := b.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		b.txSongRepo = factory.NewSongRepository()
		b.txPlaylistRepo = factory.NewPlaylistRepository()

		if err := b.reorderAlbums(songs, idsMap); err != nil {
			return err
		}
		if err := b.reorderSongsInPlaylists(songs, idsMap); err != nil {
			return err
		}
		if err := b.txSongRepo.Delete(request.IDs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	if err := b.messagePublisherService.Publish(topics.SongsDeletedTopic, songs); err != nil {
		return httperror.MessagePublisherError(err)
	}

	return nil
}

func (b BulkDeleteSongs) reorderAlbums(songs []model.Song, idsMap map[uuid.UUID]bool) error {
	albumsToReorder := make(map[uuid.UUID]model.Album)
	for _, song := range songs {
		if song.Album != nil {
			albumsToReorder[song.Album.ID] = *song.Album
		}
	}

	var albumSongsToUpdate []model.Song
	for _, album := range albumsToReorder {
		songsFound := uint(0)
		for _, albumSong := range album.Songs {
			if idsMap[albumSong.ID] {
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
		if err := b.txSongRepo.UpdateAll(&albumSongsToUpdate); err != nil {
			return err
		}
	}

	return nil
}

func (b BulkDeleteSongs) reorderSongsInPlaylists(songs []model.Song, idsMap map[uuid.UUID]bool) error {
	playlistsToReorder := make(map[uuid.UUID]model.Playlist)
	for _, song := range songs {
		for _, playlist := range song.Playlists {
			playlistsToReorder[playlist.ID] = playlist
		}
	}

	var playlistSongsToUpdate []model.PlaylistSong
	for _, playlist := range playlistsToReorder {
		songsFound := uint(0)
		for _, playlistSong := range playlist.PlaylistSongs {
			if idsMap[playlistSong.SongID] {
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
		if err := b.txPlaylistRepo.UpdateAllPlaylistSongs(&playlistSongsToUpdate); err != nil {
			return err
		}
	}

	return nil
}
