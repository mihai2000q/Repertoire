package section

import (
	"errors"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/internal/deduplicate"
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
	if len(request.PartIDs) > 0 {
		errCode := c.ensurePartsBelongToSameSong(request, request.SongID)
		if errCode != nil {
			return errCode
		}
	}

	var sectionsCount int64
	if err := c.songSectionRepository.CountAllBySong(&sectionsCount, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}

	section := model.SongSection{
		ID:                uuid.New(),
		Name:              request.Name,
		SongSectionTypeID: request.TypeID,
		Order:             uint(sectionsCount),
		SongID:            request.SongID,
		SectionParts:      c.createSectionParts(request.PartIDs),
	}
	if err := c.songSectionRepository.Create(&section); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (c CreateSongSection) ensurePartsBelongToSameSong(
	request requests.CreateSongSectionRequest,
	songID uuid.UUID,
) *httperror.ErrorCode {
	var parts []model.SongPart
	if err := c.songPartRepository.GetAllByIDs(&parts, request.PartIDs); err != nil {
		return httperror.DatabaseError(err)
	}
	partIDSet := deduplicate.Deduplicate(request.PartIDs)
	if len(parts) != len(partIDSet) {
		return httperror.NotFoundError(errors.New("some parts not found"))
	}
	for _, p := range parts {
		if p.SongID != songID {
			return httperror.ConflictError(errors.New("song part does not belong to the same song as the section"))
		}
	}
	return nil
}

func (c CreateSongSection) createSectionParts(partIDs []uuid.UUID) []model.SongSectionPart {
	var sectionParts []model.SongSectionPart
	for i, id := range partIDs {
		sectionPart := model.SongSectionPart{
			PartID: id,
			Order:  uint(i),
		}
		sectionParts = append(sectionParts, sectionPart)
	}
	return sectionParts
}
