package requests

import "github.com/google/uuid"

// Guitar Tunings

type CreateGuitarTuningRequest struct {
	Name string `validate:"required,max=16"`
}

type MoveGuitarTuningRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
}

// Song Section Types

type CreateSongSectionTypeRequest struct {
	Name string `validate:"required,max=16"`
}

type MoveSongSectionTypeRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
}
