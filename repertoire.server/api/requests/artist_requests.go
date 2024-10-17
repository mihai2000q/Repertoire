package requests

import "github.com/google/uuid"

type GetArtistsRequest struct {
	UserID      uuid.UUID `validate:"required"`
	CurrentPage *int      `validate:"gt=0"`
	PageSize    *int      `validate:"gt=0"`
}

type CreateArtistRequest struct {
	Name string `validate:"required,max=100"`
}

type UpdateArtistRequest struct {
	ID   uuid.UUID `validate:"required"`
	Name string    `validate:"required,max=100"`
}
