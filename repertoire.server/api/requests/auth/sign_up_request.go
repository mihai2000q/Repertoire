package auth

type SignUpRequest struct {
	Name     string `validate:"required,max=100"`
	Email    string `validate:"required,max=256,email"`
	Password string `validate:"required,min=8"`
}
