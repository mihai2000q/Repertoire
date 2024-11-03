package requests

import (
	"time"

	"github.com/google/uuid"
)

type GetAlbumsRequest struct {
	CurrentPage *int `validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int `validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     string
}

type CreateAlbumRequest struct {
	Title       string `validate:"required,max=100"`
	ReleaseDate *time.Time
}

type AddSongToAlbumRequest struct {
	ID     uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}

type UpdateAlbumRequest struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required,max=100"`
	ReleaseDate *time.Time
}

type MoveSongFromAlbumRequest struct {
	ID         uuid.UUID `validate:"required"`
	SongID     uuid.UUID `validate:"required"`
	OverSongID uuid.UUID `validate:"required"`
}
