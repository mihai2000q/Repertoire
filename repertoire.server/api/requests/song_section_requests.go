package requests

import "github.com/google/uuid"

type GetSongSectionsRequest struct {
	SongID uuid.UUID `form:"songId" validate:"required"`
}

type CreateSongSectionRequest struct {
	SongID uuid.UUID                      `validate:"required"`
	Name   string                         `validate:"required,max=30"`
	TypeID uuid.UUID                      `validate:"required"`
	Parts  []CreateSongSectionPartRequest `validate:"dive"`
}

type UpdateSongSectionRequest struct {
	ID      uuid.UUID `validate:"required"`
	Name    string    `validate:"required,max=30"`
	TypeID  uuid.UUID `validate:"required"`
	PartIDs []uuid.UUID
}

type MoveSongSectionRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}

type BulkDeleteSongSectionsRequest struct {
	IDs     []uuid.UUID `validate:"unique,min=1"`
	SongID  uuid.UUID   `validate:"required"`
	PartIDs []uuid.UUID
}

type CreateSongSectionPartRequest struct {
	PartID       *uuid.UUID                `validate:"excluded_with=NewPart"`
	NewPart      *CreateNewSongPartRequest `validate:"excluded_with=PartID"`
	BandMemberID *uuid.UUID                `validate:"excluded_without_all=PartID NewPart"`
}

type CreateNewSongPartRequest struct {
	Name         string `validate:"required,max=30"`
	InstrumentID *uuid.UUID
}
