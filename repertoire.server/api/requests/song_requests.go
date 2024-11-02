package requests

import (
	"repertoire/server/internal/enums"
	"time"

	"github.com/google/uuid"
)

type GetSongsRequest struct {
	CurrentPage *int `validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int `validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     string
}

type CreateSongRequest struct {
	Title          string `validate:"required,max=100"`
	Description    string
	Bpm            *uint
	SongsterrLink  *string `validate:"omitempty,url,contains=songsterr.com"`
	ReleaseDate    *time.Time
	Difficulty     *enums.Difficulty `validate:"omitempty,isDifficultyEnum"`
	GuitarTuningID *uuid.UUID
	Sections       []CreateSectionRequest `validate:"dive"`
	AlbumID        *uuid.UUID             `validate:"required_with=TrackNo,omitempty,excluded_with=AlbumTitle"`
	AlbumTitle     *string                `validate:"required_with=TrackNo,omitempty,excluded_with=AlbumID,max=100"`
	TrackNo        *uint                  `validate:"omitempty,gt=0"`
	ArtistID       *uuid.UUID             `validate:"omitempty,excluded_with=ArtistName"`
	ArtistName     *string                `validate:"omitempty,excluded_with=ArtistID,max=100"`
}

type UpdateSongRequest struct {
	ID             uuid.UUID `validate:"required"`
	Title          string    `validate:"required,max=100"`
	Description    string
	IsRecorded     bool
	Bpm            *uint
	SongsterrLink  *string `validate:"omitempty,url,contains=songsterr.com"`
	ReleaseDate    *time.Time
	Difficulty     *enums.Difficulty `validate:"omitempty,isDifficultyEnum"`
	GuitarTuningID *uuid.UUID
	AlbumID        *uuid.UUID `validate:"required_with=TrackNo"`
	TrackNo        *uint      `validate:"omitempty,gt=0"`
	ArtistID       *uuid.UUID
}

type CreateSectionRequest struct {
	Name   string    `validate:"required,max=30"`
	TypeID uuid.UUID `validate:"required"`
}

// Sections

type CreateSongSectionRequest struct {
	SongID uuid.UUID `validate:"required"`
	Name   string    `validate:"required,max=30"`
	TypeID uuid.UUID `validate:"required"`
}

type UpdateSongSectionRequest struct {
	ID         uuid.UUID `validate:"required"`
	Name       string    `validate:"required,max=30"`
	Rehearsals uint
	TypeID     uuid.UUID `validate:"required"`
}

type MoveSongSectionRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}
