package request

import "github.com/google/uuid"

type GetSongsRequest struct {
	CurrentPage *int `validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int `validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     string
}

type CreateSongRequest struct {
	Title string `validate:"required,max=100"`
}

type UpdateSongRequest struct {
	ID         uuid.UUID `validate:"required"`
	Title      string    `validate:"required,max=100"`
	IsRecorded bool
}
