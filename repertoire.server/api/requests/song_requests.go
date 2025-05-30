package requests

import (
	"github.com/google/uuid"
	"repertoire/server/internal"
	"repertoire/server/internal/enums"
)

type GetSongsRequest struct {
	CurrentPage *int     `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int     `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string `form:"orderBy" validate:"order_by"`
	SearchBy    []string `form:"searchBy" validate:"search_by"`
}

type GetSongFiltersMetadataRequest struct {
	SearchBy []string `form:"searchBy" validate:"search_by"`
}

type CreateSongRequest struct {
	Title          string `validate:"required,max=100"`
	Description    string
	Bpm            *uint
	SongsterrLink  *string `validate:"omitempty,url,contains=songsterr.com"`
	YoutubeLink    *string `validate:"omitempty,youtube_link"`
	ReleaseDate    *internal.Date
	Difficulty     *enums.Difficulty `validate:"omitempty,difficulty_enum"`
	GuitarTuningID *uuid.UUID
	Sections       []CreateSectionRequest `validate:"dive"`
	AlbumID        *uuid.UUID             `validate:"omitempty,excluded_with=AlbumTitle ArtistID ArtistName"`
	AlbumTitle     *string                `validate:"omitempty,excluded_with=AlbumID,max=100"`
	ArtistID       *uuid.UUID             `validate:"omitempty,excluded_with=ArtistName"`
	ArtistName     *string                `validate:"omitempty,excluded_with=ArtistID,max=100"`
}

type AddPerfectSongRehearsalRequest struct {
	ID uuid.UUID `validate:"required"`
}

type AddPartialSongRehearsalRequest struct {
	ID uuid.UUID `validate:"required"`
}

type UpdateSongRequest struct {
	ID             uuid.UUID `validate:"required"`
	Title          string    `validate:"required,max=100"`
	Description    string
	IsRecorded     bool
	Bpm            *uint
	SongsterrLink  *string `validate:"omitempty,url,contains=songsterr.com"`
	YoutubeLink    *string `validate:"omitempty,youtube_link"`
	ReleaseDate    *internal.Date
	Difficulty     *enums.Difficulty `validate:"omitempty,difficulty_enum"`
	GuitarTuningID *uuid.UUID
	ArtistID       *uuid.UUID
	AlbumID        *uuid.UUID
}

type UpdateSongSettingsRequest struct {
	SettingsID          uuid.UUID `validate:"required"`
	DefaultInstrumentID *uuid.UUID
	DefaultBandMemberID *uuid.UUID
}

type CreateSectionRequest struct {
	Name   string    `validate:"required,max=30"`
	TypeID uuid.UUID `validate:"required"`
}

// Sections

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

type UpdateSectionOccurrencesRequest struct {
	ID          uuid.UUID `validate:"required"`
	Occurrences uint
}

type UpdateSectionPartialOccurrencesRequest struct {
	ID                 uuid.UUID `validate:"required"`
	PartialOccurrences uint
}
