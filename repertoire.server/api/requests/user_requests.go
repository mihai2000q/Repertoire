package requests

import "github.com/google/uuid"

type UpdateUserRequest struct {
	ID   uuid.UUID `validate:"required"`
	Name string    `validate:"required"`
}
