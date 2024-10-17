package requests

import "github.com/google/uuid"

type GetSongsRequest struct {
	UserID      uuid.UUID `validate:"required"`
	CurrentPage *int      `validate:"omitempty,required_with=PageSize,gt=0"`
	PageSize    *int      `validate:"omitempty,required_with=CurrentPage,gt=0"`
}

type CreateSongRequest struct {
	Title      string `validate:"required,max=100"`
	IsRecorded *bool
}

type UpdateSongRequest struct {
	ID         uuid.UUID `validate:"required"`
	Title      string    `validate:"required,max=100"`
	IsRecorded *bool
}
