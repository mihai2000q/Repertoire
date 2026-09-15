package song

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
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

		// map for easy lookup
		idsMap := make(map[uuid.UUID]bool)
		for _, iD := range request.IDs {
			idsMap[iD] = true
		}

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
	albumsReordered := make(map[uuid.UUID]bool)
	var albumSongsToUpdate []model.Song
	for _, song := range songs {
		if song.Album == nil || albumsReordered[song.Album.ID] {
			continue
		}
		album := song.Album
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
		albumsReordered[song.Album.ID] = true
	}

	if len(albumSongsToUpdate) != 0 {
		if err := b.txSongRepo.UpdateAll(&albumSongsToUpdate); err != nil {
			return err
		}
	}

	return nil
}

func (b BulkDeleteSongs) reorderSongsInPlaylists(songs []model.Song, idsMap map[uuid.UUID]bool) error {
	playlistsReordered := make(map[uuid.UUID]bool)
	var playlistSongsToUpdate []model.PlaylistSong
	for _, song := range songs {
		for _, playlist := range song.Playlists {
			if playlistsReordered[playlist.ID] {
				continue
			}
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
			playlistsReordered[playlist.ID] = true
		}
	}

	if len(playlistSongsToUpdate) != 0 {
		if err := b.txPlaylistRepo.UpdateAllPlaylistSongs(&playlistSongsToUpdate); err != nil {
			return err
		}
	}

	return nil
}
