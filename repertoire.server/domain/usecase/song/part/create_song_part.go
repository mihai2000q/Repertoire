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

type CreateSongPart struct {
	songSectionRepository repository.SongSectionRepository
	songRepository        repository.SongRepository
	artistRepository      repository.ArtistRepository
	transactionManager    transaction.Manager

	txSongRepo            repository.SongRepository
	txSongPartRepo        repository.SongPartRepository
	txSongArrangementRepo repository.SongArrangementRepository
}

func NewCreateSongPart(
	songSectionRepository repository.SongSectionRepository,
	songRepository repository.SongRepository,
	artistRepository repository.ArtistRepository,
	transactionManager transaction.Manager,
) CreateSongPart {
	return CreateSongPart{
		songSectionRepository: songSectionRepository,
		songRepository:        songRepository,
		artistRepository:      artistRepository,
		transactionManager:    transactionManager,
	}
}

func (c CreateSongPart) Handle(request requests.CreateSongPartRequest) *httperror.ErrorCode {
	if len(request.SectionIDs) > 0 {
		if errCode := c.validateSections(request); errCode != nil {
			return errCode
		}
	}

	var song model.Song
	if err := c.songRepository.Get(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

	if request.BandMemberID != nil {
		if errCode := c.validateBandMember(*request.BandMemberID, song); errCode != nil {
			return errCode
		}
	}

	err := c.transactionManager.Execute(func(factory transaction.RepositoryFactory) error {
		c.txSongRepo = factory.NewSongRepository()
		c.txSongPartRepo = factory.NewSongPartRepository()
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
			BandMemberID: request.BandMemberID,
			InstrumentID: request.InstrumentID,
		}
		if len(request.SectionIDs) > 0 {
			if err := c.createSectionParts(request, &part); err != nil {
				return err
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

func (c CreateSongPart) validateSections(request requests.CreateSongPartRequest) *httperror.ErrorCode {
	var sections []model.SongSection
	if err := c.songSectionRepository.GetAllByIDs(&sections, request.SectionIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(sections) != len(request.SectionIDs) {
		return httperror.NotFoundError(errors.New("sections not found"))
	}
	for _, section := range sections {
		if section.SongID != request.SongID {
			return httperror.ConflictError(errors.New("section does not belong to the same song"))
		}
	}
	return nil
}

func (c CreateSongPart) validateBandMember(id uuid.UUID, song model.Song) *httperror.ErrorCode {
	var member model.BandMember
	if err := c.artistRepository.GetBandMember(&member, id); err != nil {
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

// Compute the part's order for each section: current count of parts in that section
func (c CreateSongPart) createSectionParts(request requests.CreateSongPartRequest, part *model.SongPart) error {
	counts, err := c.txSongPartRepo.CountBySectionIDs(request.SectionIDs)
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

	part.SectionParts = sectionParts

	return nil
}

// Update song's new confidence, rehearsals and progress medians
func (c CreateSongPart) updateSong(song *model.Song, part model.SongPart) error {
	songPartsCount := part.SongOrder

	song.Confidence = (song.Confidence*float64(songPartsCount) + float64(part.Confidence)) / float64(songPartsCount+1)
	song.Rehearsals = (song.Rehearsals*float64(songPartsCount) + float64(part.Rehearsals)) / float64(songPartsCount+1)
	song.Progress = (song.Progress*float64(songPartsCount) + float64(part.Progress)) / float64(songPartsCount+1)

	if err := c.txSongRepo.Update(song); err != nil {
		return err
	}

	return nil
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

	if err := c.txSongArrangementRepo.UpdateAllWithAssociations(&arrangements); err != nil {
		return err
	}

	return nil
}
