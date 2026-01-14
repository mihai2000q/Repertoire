package album

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddPerfectRehearsalsToAlbums struct {
	repository         repository.AlbumRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectRehearsalsToAlbums(
	repository repository.AlbumRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectRehearsalsToAlbums {
	return AddPerfectRehearsalsToAlbums{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectRehearsalsToAlbums) Handle(request requests.AddPerfectRehearsalsToAlbumsRequest) *wrapper.ErrorCode {
	var albums []model.Album
	err := a.repository.GetAllByIDsWithSongSections(&albums, request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(albums) == 0 {
		return wrapper.NotFoundError(errors.New("albums not found"))
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongSectionRepository := factory.NewSongSectionRepository()
		transactionSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, album := range albums {
			for _, song := range album.Songs {
				errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, transactionSongSectionRepository)
				if errC != nil {
					errCode = errC
					return errCode.Error
				}
				if isUpdated {
					newSongs = append(newSongs, song)
				}
			}
		}

		if len(newSongs) > 0 {
			err = transactionSongRepository.UpdateAllWithAssociations(&newSongs)
			if err != nil {
				errCode = wrapper.InternalServerError(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		if errCode != nil {
			return errCode
		}
		return wrapper.InternalServerError(err)
	}

	return nil
}
