package requests

import (
	"net/http"
	"repertoire/auth/api/validation"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validEmail = "someone@yahoo.com"
var validPassword = "Password123"

func TestValidateSignInRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := SignInRequest{
		Email:    validEmail,
		Password: validPassword,
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateSignInRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              SignInRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Email Test Cases
		{
			"Email is invalid because it's required",
			SignInRequest{Email: "", Password: validPassword},
			"Email",
			"required",
		},
		{
			"Email is invalid because it has too many characters",
			SignInRequest{Email: strings.Repeat("a", 257), Password: validPassword},
			"Email",
			"max",
		},
		{
			"Email is invalid because it is not an email",
			SignInRequest{Email: "someone", Password: validPassword},
			"Email",
			"email",
		},
		{
			"Email is invalid because it is not an email",
			SignInRequest{Email: "someone@yahoo", Password: validPassword},
			"Email",
			"email",
		},
		{
			"Email is invalid because it is not an email",
			SignInRequest{Email: "someone.com", Password: validPassword},
			"Email",
			"email",
		},
		// Password Test Cases
		{
			"Password is invalid because it's required",
			SignInRequest{Email: validEmail, Password: ""},
			"Password",
			"required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, errCode.Error.Error(), "SignInRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
