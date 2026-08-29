package requests

import "github.com/google/uuid"

type CreateSongSectionRequest struct {
	SongID uuid.UUID `validate:"required"`
	Name   string    `validate:"required,max=30"`
	TypeID uuid.UUID `validate:"required"`
}

type UpdateSongSectionRequest struct {
	ID     uuid.UUID `validate:"required"`
	Name   string    `validate:"required,max=30"`
	TypeID uuid.UUID `validate:"required"`
}

type MoveSongSectionRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}

type BulkDeleteSongSectionsRequest struct {
	IDs    []uuid.UUID `validate:"min=1"`
	SongID uuid.UUID   `validate:"required"`
}
