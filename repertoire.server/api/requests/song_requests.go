package requests

import "github.com/google/uuid"

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
	GuitarTuningID *uuid.UUID
	Sections       []CreateSongSectionRequest
	AlbumID        *uuid.UUID `validate:"omitempty,excluded_with=AlbumTitle"`
	AlbumTitle     *string    `validate:"omitempty,excluded_with=AlbumID,max=100"`
	ArtistID       *uuid.UUID `validate:"omitempty,excluded_with=ArtistName"`
	ArtistName     *string    `validate:"omitempty,excluded_with=ArtistID,max=100"`
}

type UpdateSongRequest struct {
	ID             uuid.UUID `validate:"required"`
	Title          string    `validate:"required,max=100"`
	Description    string
	IsRecorded     bool
	Bpm            *uint
	SongsterrLink  *string `validate:"omitempty,url,contains=songsterr.com"`
	GuitarTuningID *uuid.UUID
}

type CreateSongSectionRequest struct {
	Name   string    `validate:"required,max=30"`
	TypeId uuid.UUID `validate:"required"`
}
