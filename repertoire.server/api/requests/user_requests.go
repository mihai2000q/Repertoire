package requests

type UpdateUserRequest struct {
	Name string `validate:"required"`
}
