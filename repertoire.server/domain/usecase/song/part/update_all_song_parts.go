package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type UpdateAllSongParts struct {
	songRepository     repository.SongRepository
	artistRepository   repository.ArtistRepository
	transactionManager transaction.Manager
}

func NewUpdateAllSongParts(
	songRepository repository.SongRepository,
	artistRepository repository.ArtistRepository,
	transactionManager transaction.Manager,
) UpdateAllSongParts {
	return UpdateAllSongParts{
		songRepository:     songRepository,
		artistRepository:   artistRepository,
		transactionManager: transactionManager,
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

	if request.BandMemberID != nil {
		if errCode := u.validateBandMember(*request.BandMemberID, song); errCode != nil {
			return errCode
		}
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

		if request.BandMemberID != nil {
			var sectionParts []model.SongSectionPart
			if err := txSongSectionRepo.GetAllSectionPartsByPartIDs(&sectionParts, partIDs); err != nil {
				return err
			}

			for i := range sectionParts {
				sectionParts[i].BandMemberID = request.BandMemberID
			}

			if err := txSongSectionRepo.UpdateAllSectionParts(&sectionParts); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (u UpdateAllSongParts) validateBandMember(
	bandMemberID uuid.UUID,
	song model.Song,
) *httperror.ErrorCode {
	var member model.BandMember
	if err := u.artistRepository.GetBandMember(&member, bandMemberID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(member).IsZero() {
		return httperror.NotFoundError(errors.New("band member not found"))
	}
	if song.ArtistID == nil || *song.ArtistID != member.ArtistID {
		return httperror.ConflictError(errors.New("band member is not part of the artist associated with this song"))
	}
	return nil
}
