package artist

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/processor"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"
)

type AddPerfectRehearsalsToArtists struct {
	artistRepository   repository.ArtistRepository
	songProcessor      processor.SongProcessor
	transactionManager transaction.Manager
}

func NewAddPerfectRehearsalsToArtists(
	artistRepository repository.ArtistRepository,
	songProcessor processor.SongProcessor,
	transactionManager transaction.Manager,
) AddPerfectRehearsalsToArtists {
	return AddPerfectRehearsalsToArtists{
		artistRepository:   artistRepository,
		songProcessor:      songProcessor,
		transactionManager: transactionManager,
	}
}

func (a AddPerfectRehearsalsToArtists) Handle(request requests.AddPerfectRehearsalsToArtistsRequest) *httperror.ErrorCode {
	var artists []model.Artist
	if err := a.artistRepository.GetAllByIDsWithSongPartsAndDefaultOccurrences(&artists, request.IDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(artists) != len(request.IDs) {
		return httperror.NotFoundError(errors.New("artists not found"))
	}

	var errCode *httperror.ErrorCode
	err := a.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongPartRepo := factory.NewSongPartRepository()
		txSongRepo := factory.NewSongRepository()

		var newSongs []model.Song
		for _, artist := range artists {
			for _, song := range artist.Songs {
				errC, isUpdated := a.songProcessor.AddPerfectRehearsal(&song, txSongPartRepo)
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
			if err := txSongRepo.UpdateAllWithAssociations(&newSongs); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		if errCode != nil {
			return errCode
		}
		return httperror.DatabaseError(err)
	}

	return nil
}
