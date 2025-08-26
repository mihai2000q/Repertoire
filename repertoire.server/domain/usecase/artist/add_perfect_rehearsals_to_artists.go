package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
)

type AddPerfectRehearsalsToArtists struct {
	repository         repository.ArtistRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectRehearsalsToArtists(
	repository repository.ArtistRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectRehearsalsToArtists {
	return AddPerfectRehearsalsToArtists{
		repository:         repository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectRehearsalsToArtists) Handle(request requests.AddPerfectRehearsalsToArtistsRequest) *wrapper.ErrorCode {
	var artists []model.Artist
	err := a.repository.GetAllByIDsWithSongSections(&artists, request.IDs)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if len(artists) == 0 {
		return wrapper.NotFoundError(errors.New("artists not found"))
	}

	var errCode *wrapper.ErrorCode
	err = a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		transactionSongRepository := factory.NewSongRepository()

		var newSongs []model.Song
		for _, artist := range artists {
			for _, song := range artist.Songs {
				errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, transactionSongRepository)
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
