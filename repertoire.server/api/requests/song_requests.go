package requests

import "github.com/google/uuid"

type GetSongsRequest struct {
	CurrentPage *int `validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int `validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     string
}

type CreateSongSectionRequest struct {
	Name   string    `validate:"required,max=30"`
	TypeId uuid.UUID `validate:"required"`
}

type CreateSongRequest struct {
	Title          string `validate:"required,max=100"`
	Description    string
	Bpm            *uint
	SongsterrLink  *string `validate:"omitempty,url,contains=songsterr.com"`
	GuitarTuningID *uuid.UUID
}

type UpdateSongRequest struct {
	ID             uuid.UUID `validate:"required"`
	Title          string    `validate:"required,max=100"`
	Description    string
	IsRecorded     bool
	Rehearsals     uint
	Bpm            *uint
	SongsterrLink  *string `validate:"omitempty,url,contains=songsterr.com"`
	GuitarTuningID *uuid.UUID
}
