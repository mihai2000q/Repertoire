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

type CreateSongPart struct {
	songSectionRepository repository.SongSectionRepository
	songRepository        repository.SongRepository
	bandMemberValidator   validator.BandMemberValidator
	transactionManager    transaction.Manager

	txSongRepo            repository.SongRepository
	txSongPartRepo        repository.SongPartRepository
	txSongSectionRepo     repository.SongSectionRepository
	txSongArrangementRepo repository.SongArrangementRepository
}

func NewCreateSongPart(
	songSectionRepository repository.SongSectionRepository,
	songRepository repository.SongRepository,
	bandMemberValidator validator.BandMemberValidator,
	transactionManager transaction.Manager,
) CreateSongPart {
	return CreateSongPart{
		songSectionRepository: songSectionRepository,
		songRepository:        songRepository,
		bandMemberValidator:   bandMemberValidator,
		transactionManager:    transactionManager,
	}
}

func (c CreateSongPart) Handle(request requests.CreateSongPartRequest) *httperror.ErrorCode {
	var song model.Song
	if err := c.songRepository.Get(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	if request.SectionID != nil {
		if errCode := c.validateSection(*request.SectionID, song); errCode != nil {
			return errCode
		}
	}

	bandMembers, errCode := c.bandMemberValidator.Validate(request.BandMemberIDs, song)
	if errCode != nil {
		return errCode
	}

	err := c.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		c.txSongRepo = factory.NewSongRepository()
		c.txSongPartRepo = factory.NewSongPartRepository()
		c.txSongSectionRepo = factory.NewSongSectionRepository()
		c.txSongArrangementRepo = factory.NewSongArrangementRepository()

		var songPartsCount int64
		if err := c.txSongPartRepo.CountAllBySong(&songPartsCount, request.SongID); err != nil {
			return err
		}

		part := model.SongPart{
			ID:           uuid.New(),
			Name:         request.Name,
			Confidence:   model.DefaultSongPartConfidence,
			SongOrder:    uint(songPartsCount),
			SongID:       request.SongID,
			InstrumentID: request.InstrumentID,
		}

		if request.SectionID != nil {
			var sectionPartsCount int64
			if err := c.txSongSectionRepo.CountAllSectionPartsBySectionID(&sectionPartsCount, *request.SectionID); err != nil {
				return err
			}

			part.SectionParts = []model.SongSectionPart{
				{
					PartID:      part.ID,
					SectionID:   *request.SectionID,
					Order:       uint(sectionPartsCount),
					BandMembers: bandMembers,
				},
			}
		}

		if err := c.txSongPartRepo.Create(&part); err != nil {
			return err
		}

		if err := c.updateSong(&song, part); err != nil {
			return err
		}

		if err := c.updateArrangements(part.ID, request.SongID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (c CreateSongPart) validateSection(sectionID uuid.UUID, song model.Song) *httperror.ErrorCode {
	var section model.SongSection
	if err := c.songSectionRepository.Get(&section, sectionID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(section).IsZero() {
		return httperror.NotFoundError(errors.New("section not found"))
	}
	if section.SongID != song.ID {
		return httperror.ConflictError(errors.New("section does not belong to the same song"))
	}
	return nil
}

// Update song's new confidence, rehearsals and progress medians
func (c CreateSongPart) updateSong(song *model.Song, part model.SongPart) error {
	songPartsCount := part.SongOrder

	song.Confidence = (song.Confidence*float64(songPartsCount) + float64(part.Confidence)) / float64(songPartsCount+1)
	song.Rehearsals = (song.Rehearsals*float64(songPartsCount) + float64(part.Rehearsals)) / float64(songPartsCount+1)
	song.Progress = (song.Progress*float64(songPartsCount) + float64(part.Progress)) / float64(songPartsCount+1)

	return c.txSongRepo.Update(song)
}

// Add one new part occurrence on each song arrangement
func (c CreateSongPart) updateArrangements(partID uuid.UUID, songID uuid.UUID) error {
	var arrangements []model.SongArrangement
	if err := c.txSongArrangementRepo.GetAllBySong(&arrangements, songID); err != nil {
		return err
	}

	for i := range arrangements {
		occurrence := model.SongPartOccurrences{
			PartID:        partID,
			Occurrences:   0,
			ArrangementID: arrangements[i].ID,
		}
		arrangements[i].PartOccurrences = append(arrangements[i].PartOccurrences, occurrence)
	}

	return c.txSongArrangementRepo.UpdateAllWithAssociations(&arrangements)
}
