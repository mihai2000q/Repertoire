package part

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/database/transaction"
	"repertoire/server/data/repository"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongPart struct {
	songSectionRepository repository.SongSectionRepository
	songRepository        repository.SongRepository
	transactionManager    transaction.Manager

	txSongRepository            repository.SongRepository
	txSongSectionRepository     repository.SongSectionRepository
	txSongPartRepository        repository.SongPartRepository
	txSongArrangementRepository repository.SongArrangementRepository
}

func NewCreateSongPart(
	songSectionRepository repository.SongSectionRepository,
	songRepository repository.SongRepository,
	transactionManager transaction.Manager,
) CreateSongPart {
	return CreateSongPart{
		songSectionRepository: songSectionRepository,
		songRepository:        songRepository,
		transactionManager:    transactionManager,
	}
}

func (c CreateSongPart) Handle(request requests.CreateSongPartRequest) *wrapper.ErrorCode {
	// Validate sections
	if len(request.SectionIDs) > 0 {
		var sections []model.SongSection
		if err := c.songSectionRepository.GetAllByIDs(&sections, request.SectionIDs); err != nil {
			return wrapper.InternalServerError(err)
		}
		if len(sections) != len(request.SectionIDs) {
			return wrapper.NotFoundError(errors.New("sections not found"))
		}
		for _, section := range sections {
			if section.SongID != request.SongID {
				return wrapper.ConflictError(errors.New("section does not belong to the same song"))
			}
		}
	}

	// Validate band member
	if request.BandMemberID != nil {
		res, err := c.songRepository.IsBandMemberAssociatedWithSong(request.SongID, *request.BandMemberID)
		if err != nil {
			return wrapper.InternalServerError(err)
		}
		if !res {
			return wrapper.ConflictError(errors.New("band member is not part of the artist associated with this song"))
		}
	}

	var errCode *wrapper.ErrorCode
	err := c.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		c.txSongRepository = factory.NewSongRepository()
		c.txSongPartRepository = factory.NewSongPartRepository()
		c.txSongArrangementRepository = factory.NewSongArrangementRepository()

		var songPartsCount int64
		err := c.txSongPartRepository.CountAllBySong(&songPartsCount, request.SongID)
		if err != nil {
			return err
		}

		part := model.SongPart{
			ID:           uuid.New(),
			Name:         request.Name,
			Confidence:   model.DefaultSongPartConfidence,
			SongOrder:    uint(songPartsCount),
			SongID:       request.SongID,
			BandMemberID: request.BandMemberID,
			InstrumentID: request.InstrumentID,
		}
		if err = c.txSongPartRepository.Create(&part); err != nil {
			return err
		}

		if len(request.SectionIDs) > 0 {
			c.txSongSectionRepository = factory.NewSongSectionRepository()
			if err = c.createSectionParts(request, part); err != nil {
				return err
			}
		}

		errCode = c.updateSong(part)
		if errCode != nil {
			return errCode.Error
		}

		errCode = c.updateArrangements(part.ID, request.SongID)
		if errCode != nil {
			return errCode.Error
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

// Compute the part's order for each section: current count of parts in that section
func (c CreateSongPart) createSectionParts(request requests.CreateSongPartRequest, part model.SongPart) error {
	counts, err := c.txSongPartRepository.CountBySectionIDs(request.SectionIDs)
	if err != nil {
		return err
	}

	sectionParts := make([]model.SongSectionPart, len(request.SectionIDs))
	for i, sectionID := range request.SectionIDs {
		order := uint(counts[sectionID])
		sectionParts[i] = model.SongSectionPart{
			PartID:    part.ID,
			SectionID: sectionID,
			Order:     order,
		}
	}

	if err = c.txSongSectionRepository.CreateAllSectionParts(&sectionParts); err != nil {
		return err
	}

	return nil
}

// Update song's new confidence, rehearsals and progress medians
func (c CreateSongPart) updateSong(part model.SongPart) *wrapper.ErrorCode {
	var song model.Song
	err := c.txSongRepository.Get(&song, part.SongID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return wrapper.NotFoundError(errors.New("song not found"))
	}

	songPartsCount := part.SongOrder

	song.Confidence = (song.Confidence*float64(songPartsCount) + float64(part.Confidence)) / float64(songPartsCount+1)
	song.Rehearsals = (song.Rehearsals*float64(songPartsCount) + float64(part.Rehearsals)) / float64(songPartsCount+1)
	song.Progress = (song.Progress*float64(songPartsCount) + float64(part.Progress)) / float64(songPartsCount+1)

	err = c.txSongRepository.Update(&song)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}

// Add one new part occurrence on each song arrangement
func (c CreateSongPart) updateArrangements(partID uuid.UUID, songID uuid.UUID) *wrapper.ErrorCode {
	var arrangements []model.SongArrangement
	err := c.txSongArrangementRepository.GetAllBySong(&arrangements, songID)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	for i := range arrangements {
		occurrence := model.SongPartOccurrences{
			PartID:        partID,
			Occurrences:   0,
			ArrangementID: arrangements[i].ID,
		}
		arrangements[i].PartOccurrences = append(arrangements[i].PartOccurrences, occurrence)
	}

	err = c.txSongArrangementRepository.UpdateAllWithAssociations(&arrangements)
	if err != nil {
		return wrapper.InternalServerError(err)
	}

	return nil
}
