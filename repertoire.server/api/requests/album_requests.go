package requests

import "github.com/google/uuid"

type GetAlbumsRequest struct {
	UserID      uuid.UUID `validate:"required"`
	CurrentPage *int      `validate:"omitempty,required_with=PageSize,gt=0"`
	PageSize    *int      `validate:"omitempty,required_with=CurrentPage,gt=0"`
}

type CreateAlbumRequest struct {
	Title string `validate:"required,max=100"`
}

type UpdateAlbumRequest struct {
	ID    uuid.UUID `validate:"required"`
	Title string    `validate:"required,max=100"`
}
