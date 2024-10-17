package requests

import "github.com/google/uuid"

type GetPlaylistsRequest struct {
	UserID      uuid.UUID `validate:"required"`
	CurrentPage *int      `validate:"omitempty,required_with=PageSize,gt=0"`
	PageSize    *int      `validate:"omitempty,required_with=CurrentPage,gt=0"`
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
