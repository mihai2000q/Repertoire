package section

import (
	"errors"
	"reflect"
	"repertoire/server/api/requests"
	"repertoire/server/data/repository"
	"repertoire/server/domain/validator"
	"repertoire/server/internal/httperror"
	"repertoire/server/model"

	"github.com/google/uuid"
)

type CreateSongSection struct {
	songSectionRepository repository.SongSectionRepository
	songPartRepository    repository.SongPartRepository
	songRepository        repository.SongRepository
	bandMemberValidator   validator.BandMemberValidator
	songPartValidator     validator.SongPartValidator
}

func NewCreateSongSection(
	songSectionRepository repository.SongSectionRepository,
	songPartRepository repository.SongPartRepository,
	songRepository repository.SongRepository,
	bandMemberValidator validator.BandMemberValidator,
	songPartValidator validator.SongPartValidator,
) CreateSongSection {
	return CreateSongSection{
		songSectionRepository: songSectionRepository,
		songPartRepository:    songPartRepository,
		songRepository:        songRepository,
		bandMemberValidator:   bandMemberValidator,
		songPartValidator:     songPartValidator,
	}
}

func (c CreateSongSection) Handle(request requests.CreateSongSectionRequest) *httperror.ErrorCode {
	var song model.Song
	if err := c.songRepository.Get(&song, request.SongID); err != nil {
		return httperror.DatabaseError(err)
	}
	if reflect.ValueOf(song).IsZero() {
		return httperror.NotFoundError(errors.New("song not found"))
	}

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
		if errCode := c.songPartValidator.Validate(partIDs, request.SongID); errCode != nil {
			return errCode
		}
	}

	bandMembersByID, errCode := c.validateBandMembers(request.Parts, song)
	if errCode != nil {
		return errCode
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
		SectionParts:      c.createSectionParts(request.Parts, request.SongID, uint(partsCount), bandMembersByID),
	}
	if err := c.songSectionRepository.Create(&section); err != nil {
		return httperror.DatabaseError(err)
	}

	return nil
}

func (c CreateSongSection) validateBandMembers(
	parts []requests.CreateSongSectionPartRequest,
	song model.Song,
) (map[uuid.UUID]model.BandMember, *httperror.ErrorCode) {
	seen := make(map[uuid.UUID]bool)
	ids := make([]uuid.UUID, 0)
	for _, p := range parts {
		for _, id := range p.BandMemberIDs {
			if !seen[id] {
				seen[id] = true
				ids = append(ids, id)
			}
		}
	}

	bandMembers, errCode := c.bandMemberValidator.Validate(ids, song)
	if errCode != nil {
		return nil, errCode
	}

	bandMembersMap := make(map[uuid.UUID]model.BandMember, len(bandMembers))
	for _, bandMember := range bandMembers {
		bandMembersMap[bandMember.ID] = bandMember
	}
	return bandMembersMap, nil
}

func (c CreateSongSection) createSectionParts(
	parts []requests.CreateSongSectionPartRequest,
	songID uuid.UUID,
	existingPartsCount uint,
	bandMembersMap map[uuid.UUID]model.BandMember,
) []model.SongSectionPart {
	sectionParts := make([]model.SongSectionPart, len(parts))
	nextSongOrder := existingPartsCount

	for i, p := range parts {
		sectionPart := model.SongSectionPart{
			Order:       uint(i),
			BandMembers: bandMembersForPart(p.BandMemberIDs, bandMembersMap),
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

func bandMembersForPart(ids []uuid.UUID, bandMembersByID map[uuid.UUID]model.BandMember) []model.BandMember {
	if len(ids) == 0 {
		return nil
	}
	bandMembers := make([]model.BandMember, len(ids))
	for i, id := range ids {
		bandMembers[i] = bandMembersByID[id]
	}
	return bandMembers
}
