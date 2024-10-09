package auth

type SignInRequest struct {
	Email    string `validate:"required,max=256,email"`
	Password string `validate:"required,min=8"`
}
