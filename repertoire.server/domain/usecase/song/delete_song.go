package song

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"reflect"
	"repertoire/server/data/repository"
	"repertoire/server/data/service"
	"repertoire/server/domain/provider"
	"repertoire/server/internal/message/topics"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type DeleteSong struct {
	repository              repository.SongRepository
	playlistRepository      repository.PlaylistRepository
	storageService          service.StorageService
	storageFilePathProvider provider.StorageFilePathProvider
	messagePublisherService service.MessagePublisherService
}

func NewDeleteSong(
	repository repository.SongRepository,
	playlistRepository repository.PlaylistRepository,
	storageService service.StorageService,
	storageFilePathProvider provider.StorageFilePathProvider,
	messagePublisherService service.MessagePublisherService,
) DeleteSong {
	return DeleteSong{
		repository:              repository,
		playlistRepository:      playlistRepository,
		storageService:          storageService,
		storageFilePathProvider: storageFilePathProvider,
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

	// TODO: RUN IN PARALLEL (GOROUTINES)
	if song.AlbumID != nil {
		errCode := d.reorderAlbum(song)
		if errCode != nil {
			return errCode
		}
	}
	
	directoryPath := d.storageFilePathProvider.GetSongDirectoryPath(song)
	errCode := d.storageService.DeleteDirectory(directoryPath)
	if errCode != nil && errCode.Code != http.StatusNotFound {
		return errCode
	}

	errCode = d.reorderSongsInPlaylists(song)
	if errCode != nil {
		return errCode
	}

	err = d.repository.Delete(id)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	err = d.messagePublisherService.Publish(topics.SongDeletedTopic, song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

func (d DeleteSong) reorderAlbum(song model.Song) *wrapper.ErrorCode {
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
		songFound := false
		var playlistSongsToUpdate []model.PlaylistSong
		for _, playlistSong := range playlist.PlaylistSongs {
			if songFound {
				playlistSong.SongTrackNo = playlistSong.SongTrackNo - 1
				playlistSongsToUpdate = append(playlistSongsToUpdate, playlistSong)
			}

			if playlistSong.SongID == song.ID {
				songFound = true
			}
		}
		err := d.playlistRepository.UpdateAllPlaylistSongs(&playlistSongsToUpdate)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
	}
	return nil
}
