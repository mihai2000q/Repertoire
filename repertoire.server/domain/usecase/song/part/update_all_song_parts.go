package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/domain/validator"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UpdateAllSongParts struct {
	songRepository      repository.SongRepository
	bandMemberValidator validator.BandMemberValidator
	transactionManager  transaction.Manager
}

func NewUpdateAllSongParts(
	songRepository repository.SongRepository,
	bandMemberValidator validator.BandMemberValidator,
	transactionManager transaction.Manager,
) UpdateAllSongParts {
	return UpdateAllSongParts{
		songRepository:      songRepository,
		bandMemberValidator: bandMemberValidator,
		transactionManager:  transactionManager,
	}
}

func (u UpdateAllSongParts) Handle(request requests.UpdateAllSongPartsRequest) *httperror.ErrorCode {
	var song model.Song
	if err := u.songRepository.GetWithParts(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	var bandMemberIDs []uuid.UUID
	if request.BandMemberID != nil {
		bandMemberIDs = []uuid.UUID{*request.BandMemberID}
	}
	bandMembers, errCode := u.bandMemberValidator.Validate(bandMemberIDs, song)
	if errCode != nil {
		return errCode
	}

	err := u.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		txSongRepo := factory.NewSongRepository()
		txSongSectionRepo := factory.NewSongSectionRepository()

		partIDs := make([]uuid.UUID, len(song.Parts))
		for i := range song.Parts {
			partIDs[i] = song.Parts[i].ID
		}

		if request.InstrumentID != nil {
			for i := range song.Parts {
				song.Parts[i].InstrumentID = request.InstrumentID
			}
			if err := txSongRepo.UpdateWithAssociations(&song); err != nil {
				return err
			}
		}

		// no band members in the request means leaving the current ones untouched
		if len(bandMembers) > 0 {
			var sectionParts []model.SongSectionPart
			if err := txSongSectionRepo.GetAllSectionPartsByPartIDs(&sectionParts, partIDs); err != nil {
				return err
			}

			for i := range sectionParts {
				if err := txSongSectionRepo.ReplaceSectionPartBandMembers(&sectionParts[i], bandMembers); err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}
