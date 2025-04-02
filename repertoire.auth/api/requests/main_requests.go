package requests

type RefreshRequest struct {
	Token string
}

type SignInRequest struct {
	Email    string `validate:"required,max=256,email"`
	Password string `validate:"required"`
}
