package requests

import "github.com/google/uuid"

type GetPlaylistsRequest struct {
	CurrentPage *int `validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int `validate:"required_with=CurrentPage,omitempty,gt=0"`
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
