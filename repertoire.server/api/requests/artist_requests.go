package requests

import "github.com/google/uuid"

type GetArtistsRequest struct {
	UserID uuid.UUID `validate:"required"`
}

type CreateArtistRequest struct {
	Name string `validate:"required,max=100"`
}

type UpdateArtistRequest struct {
	ID   uuid.UUID `validate:"required"`
	Name string    `validate:"required,max=100"`
}
