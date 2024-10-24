package request

import "github.com/google/uuid"

type GetArtistsRequest struct {
	CurrentPage *int `validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int `validate:"required_with=CurrentPage,omitempty,gt=0"`
	OrderBy     string
}

type CreateArtistRequest struct {
	Name string `validate:"required,max=100"`
}

type UpdateArtistRequest struct {
	ID   uuid.UUID `validate:"required"`
	Name string    `validate:"required,max=100"`
}
