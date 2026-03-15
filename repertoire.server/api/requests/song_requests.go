package requests

import (
	"repertoire/server/internal"
	"repertoire/server/internal/enums"

	"github.com/google/uuid"
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

type AddPerfectSongRehearsalsRequest struct {
	IDs []uuid.UUID `validate:"min=1"`
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

type BulkDeleteSongsRequest struct {
	IDs []uuid.UUID `validate:"min=1"`
}

type CreateSectionRequest struct {
	Name   string    `validate:"required,max=30"`
	TypeID uuid.UUID `validate:"required"`
}
