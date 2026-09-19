package section

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongSection struct {
	songSectionRepository repository.SongSectionRepository
	songPartRepository    repository.SongPartRepository
}

func NewCreateSongSection(
	songSectionRepository repository.SongSectionRepository,
	songPartRepository repository.SongPartRepository,
) CreateSongSection {
	return CreateSongSection{
		songSectionRepository: songSectionRepository,
		songPartRepository:    songPartRepository,
	}
}

func (c CreateSongSection) Handle(request requests.CreateSongSectionRequest) *httperror.ErrorCode {
	partIDs := make([]uuid.UUID, 0, len(request.Parts))
	hasNewParts := false
	for _, p := range request.Parts {
		if p.PartID != nil {
			partIDs = append(partIDs, *p.PartID)
		} else {
			hasNewParts = true
		}
	}
	if len(partIDs) > 0 {
		if errCode := c.ensurePartsBelongToSameSong(partIDs, request.SongID); errCode != nil {
			return errCode
		}
	}

	var sectionsCount int64
	if err := c.songSectionRepository.CountAllBySong(&sectionsCount, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}

	var partsCount int64
	if hasNewParts {
		if err := c.songPartRepository.CountAllBySong(&partsCount, request.SongID); err != nil {
			return httperror.DatabaseError(err)
		}
	}

	section := model.SongSection{
		ID:                uuid.New(),
		Name:              request.Name,
		SongSectionTypeID: request.TypeID,
		Order:             uint(sectionsCount),
		SongID:            request.SongID,
		SectionParts:      c.createSectionParts(request.Parts, request.SongID, uint(partsCount)),
	}
	if err := c.songSectionRepository.Create(&section); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (c CreateSongSection) ensurePartsBelongToSameSong(
	partIDs []uuid.UUID,
	songID uuid.UUID,
) *httperror.ErrorCode {
	var parts []model.SongPart
	if err := c.songPartRepository.GetAllByIDs(&parts, partIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	if len(parts) != len(partIDs) {
		return httperror.NotFoundError(errors.New("parts not found"))
	}
	for _, p := range parts {
		if p.SongID != songID {
			return httperror.ConflictError(errors.New("song part does not belong to the same song as the section"))
		}
	}
	return nil
}

func (c CreateSongSection) createSectionParts(
	parts []requests.CreateSongSectionPartRequest,
	songID uuid.UUID,
	existingPartsCount uint,
) []model.SongSectionPart {
	sectionParts := make([]model.SongSectionPart, len(parts))
	nextSongOrder := existingPartsCount

	for i, p := range parts {
		sectionPart := model.SongSectionPart{
			Order:        uint(i),
			BandMemberID: p.BandMemberID,
		}

		if p.PartID != nil {
			sectionPart.PartID = *p.PartID
		} else {
			partID := uuid.New()
			sectionPart.PartID = partID
			sectionPart.Part = model.SongPart{
				ID:           partID,
				Name:         p.NewPart.Name,
				SongOrder:    nextSongOrder,
				SongID:       songID,
				InstrumentID: p.NewPart.InstrumentID,
			}
			nextSongOrder++
		}

		sectionParts[i] = sectionPart
	}

	return sectionParts
}
