package requests

import "github.com/google/uuid"

type GetSongsRequest struct {
	UserId uuid.UUID `validate:"required"`
}

type CreateSongRequest struct {
	Title      string `validate:"required,max=100"`
	IsRecorded *bool
}

type UpdateSongRequest struct {
	Id         uuid.UUID `validate:"required"`
	Title      string    `validate:"required,max=100"`
	IsRecorded *bool
}
