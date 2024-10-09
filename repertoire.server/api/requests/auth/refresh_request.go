package auth

type RefreshRequest struct {
	AccessToken string `validate:"required,jwt"`
}
