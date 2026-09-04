package requests

import "github.com/google/uuid"

type CreateSongPartRequest struct {
	SongID       uuid.UUID `validate:"required"`
	Name         string    `validate:"required,max=30"`
	SectionIDs   []uuid.UUID
	BandMemberID *uuid.UUID
	InstrumentID *uuid.UUID
}

type UpdateSongPartRequest struct {
	ID           uuid.UUID `validate:"required"`
	Name         string    `validate:"required,max=30"`
	Confidence   uint      `validate:"max=100"`
	Rehearsals   uint
	BandMemberID *uuid.UUID
	InstrumentID *uuid.UUID
}

type UpdateAllSongPartsRequest struct {
	SongID       uuid.UUID `validate:"required"`
	InstrumentID *uuid.UUID
	BandMemberID *uuid.UUID
}

type MoveSongPartInSongRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}

type MoveSongPartInSectionRequest struct {
	ID        uuid.UUID `validate:"required"`
	OverID    uuid.UUID `validate:"required"`
	SectionID uuid.UUID `validate:"required"`
}

type BulkUpdateSongPartsRequest struct {
	Requests []BulkUpdateSongPartRequest `validate:"min=1,dive"`
	SongID   uuid.UUID                   `validate:"required"`
}

type BulkDeleteSongPartsRequest struct {
	IDs    []uuid.UUID `validate:"min=1"`
	SongID uuid.UUID   `validate:"required"`
}

type BulkUpdateSongPartRequest struct {
	ID         uuid.UUID `validate:"required"`
	Confidence uint      `validate:"max=100"`
	Rehearsals uint
}
