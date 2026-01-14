package requests

import "github.com/google/uuid"

type CreateSongSectionRequest struct {
	SongID       uuid.UUID `validate:"required"`
	Name         string    `validate:"required,max=30"`
	TypeID       uuid.UUID `validate:"required"`
	BandMemberID *uuid.UUID
	InstrumentID *uuid.UUID
}

type UpdateSongSectionRequest struct {
	ID           uuid.UUID `validate:"required"`
	Name         string    `validate:"required,max=30"`
	Confidence   uint      `validate:"max=100"`
	Rehearsals   uint
	TypeID       uuid.UUID `validate:"required"`
	BandMemberID *uuid.UUID
	InstrumentID *uuid.UUID
}

type UpdateSongSectionsOccurrencesRequest struct {
	SongID   uuid.UUID                         `validate:"required"`
	Sections []UpdateSectionOccurrencesRequest `validate:"min=1,dive"`
}

type UpdateSongSectionsPartialOccurrencesRequest struct {
	SongID   uuid.UUID                                `validate:"required"`
	Sections []UpdateSectionPartialOccurrencesRequest `validate:"min=1,dive"`
}

type UpdateAllSongSectionsRequest struct {
	SongID       uuid.UUID `validate:"required"`
	InstrumentID *uuid.UUID
	BandMemberID *uuid.UUID
}

type MoveSongSectionRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}

type BulkRehearsalsSongSectionsRequest struct {
	Sections []BulkRehearsalsSongSectionRequest `validate:"min=1,dive"`
	SongID   uuid.UUID                          `validate:"required"`
}

type BulkDeleteSongSectionsRequest struct {
	IDs    []uuid.UUID `validate:"min=1"`
	SongID uuid.UUID   `validate:"required"`
}

type UpdateSectionOccurrencesRequest struct {
	ID          uuid.UUID `validate:"required"`
	Occurrences uint
}

type UpdateSectionPartialOccurrencesRequest struct {
	ID                 uuid.UUID `validate:"required"`
	PartialOccurrences uint
}

type BulkRehearsalsSongSectionRequest struct {
	ID         uuid.UUID `validate:"required"`
	Rehearsals uint
}
