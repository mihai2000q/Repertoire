package requests

import "github.com/google/uuid"

type GetPlaylistsRequest struct {
	CurrentPage *int     `form:"currentPage" validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int     `form:"pageSize" validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     []string `form:"orderBy"`
	SearchBy    []string `form:"searchBy"`
}

type CreatePlaylistRequest struct {
	Title       string `validate:"required,max=100"`
	Description string
}

type AddSongToPlaylistRequest struct {
	ID     uuid.UUID `validate:"required"`
	SongID uuid.UUID `validate:"required"`
}

type UpdatePlaylistRequest struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required,max=100"`
	Description string
}
