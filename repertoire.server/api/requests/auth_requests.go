package requests

type RefreshRequest struct {
	Token string
}

type SignInRequest struct {
	Email    string `validate:"required,max=256,email"`
	Password string `validate:"required"`
}

type SignUpRequest struct {
	Name     string `validate:"required,max=100"`
	Email    string `validate:"required,max=256,email"`
	Password string `validate:"required,min=8,has_upper,has_lower,has_digit"`
}
