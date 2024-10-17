package requests

import "github.com/google/uuid"

type GetArtistsRequest struct {
	UserID      uuid.UUID `validate:"required"`
	CurrentPage *int      `validate:"omitempty,required_with=PageSize,gt=0"`
	PageSize    *int      `validate:"omitempty,required_with=CurrentPage,gt=0"`
}

type CreateArtistRequest struct {
	Name string `validate:"required,max=100"`
}

type UpdateArtistRequest struct {
	ID   uuid.UUID `validate:"required"`
	Name string    `validate:"required,max=100"`
}
