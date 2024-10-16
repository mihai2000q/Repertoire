package requests

import "github.com/google/uuid"

type GetPlaylistsRequest struct {
	UserID uuid.UUID `validate:"required"`
}

type CreatePlaylistRequest struct {
	Title       string `validate:"required,max=100"`
	Description string
}

type UpdatePlaylistRequest struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required,max=100"`
	Description string
}
