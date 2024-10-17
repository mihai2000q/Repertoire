package requests

import "github.com/google/uuid"

type GetAlbumsRequest struct {
	UserID      uuid.UUID `validate:"required"`
	CurrentPage *int      `validate:"required_with=PageSize,omitempty,gt=0"`
	PageSize    *int      `validate:"required_with=CurrentPage,omitempty,gt=0"`
}

type CreateAlbumRequest struct {
	Title string `validate:"required,max=100"`
}

type UpdateAlbumRequest struct {
	ID    uuid.UUID `validate:"required"`
	Title string    `validate:"required,max=100"`
}
