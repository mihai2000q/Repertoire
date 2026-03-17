package requests

import "github.com/google/uuid"

type GetSongArrangementsRequest struct {
	SongID uuid.UUID `form:"songId" validate:"required"`
}

type CreateSongArrangementRequest struct {
	SongID uuid.UUID `validate:"required"`
	Name   string    `validate:"required,max=30"`
}

type BulkUpdateSongArrangementsRequest struct {
	Requests []UpdateSongArrangementRequest `validate:"min=1,dive"`
	SongID   uuid.UUID                      `validate:"required"`
}

type UpdateSongArrangementRequest struct {
	ID          uuid.UUID                             `validate:"required"`
	Name        string                                `validate:"required,max=30"`
	Occurrences []UpdateSongSectionOccurrencesRequest `validate:"omitempty,dive"`
}

type UpdateSongSectionOccurrencesRequest struct {
	SectionID   uuid.UUID `validate:"required"`
	Occurrences uint
}

type UpdateDefaultSongArrangementRequest struct {
	ID     *uuid.UUID
	SongID uuid.UUID `validate:"required"`
}

type MoveSongArrangementRequest struct {
	ID     uuid.UUID `validate:"required"`
	OverID uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}
