package requests

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/validation"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validName = "Samuel"

func TestValidateSignUpRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.SignUpRequest{
		Name:     validName,
		Email:    validEmail,
		Password: validPassword,
	}

	// when
	errCode := _uut.Validate(request)

	// then
	assert.Nil(t, errCode)
}

func TestValidateSignUpRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.SignUpRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.SignUpRequest{Name: "", Email: validEmail, Password: validPassword},
			"Name",
			"required",
		},
		// Email Test Cases
		{
			"Email is invalid because it's required",
			requests.SignUpRequest{Name: validName, Email: "", Password: validPassword},
			"Email",
			"required",
		},
		{
			"Email is invalid because it has too many characters",
			requests.SignUpRequest{Name: validName, Email: strings.Repeat("a", 257), Password: validPassword},
			"Email",
			"max",
		},
		{
			"Email is invalid because it is not an email",
			requests.SignUpRequest{Name: validName, Email: "someone", Password: validPassword},
			"Email",
			"email",
		},
		{
			"Email is invalid because it is not an email",
			requests.SignUpRequest{Name: validName, Email: "someone@yahoo", Password: validPassword},
			"Email",
			"email",
		},
		{
			"Email is invalid because it is not an email",
			requests.SignUpRequest{Name: validName, Email: "someone.com", Password: validPassword},
			"Email",
			"email",
		},
		// Password Test Cases
		{
			"Password is invalid because it's required",
			requests.SignUpRequest{Name: validName, Email: validEmail, Password: ""},
			"Password",
			"required",
		},
		{
			"Password is invalid because it has less than 8 characters",
			requests.SignUpRequest{Name: validName, Email: validEmail, Password: "1234567"},
			"Password",
			"min",
		},
		{
			"Password is invalid because it doesn't have an uppercase letter",
			requests.SignUpRequest{Name: validName, Email: validEmail, Password: strings.Repeat("a", 9)},
			"Password",
			"has_upper",
		},
		{
			"Password is invalid because it doesn't have a lowercase letter",
			requests.SignUpRequest{Name: validName, Email: validEmail, Password: strings.Repeat("A", 9)},
			"Password",
			"has_lower",
		},
		{
			"Password is invalid because it doesn't have any digit",
			requests.SignUpRequest{Name: validName, Email: validEmail, Password: strings.Repeat("A", 4) + strings.Repeat("a", 4)},
			"Password",
			"has_digit",
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
			assert.Contains(t, errCode.Error.Error(), "SignUpRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}

func TestValidateUpdateUserRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateUserRequest{
		Name: validName,
	}

	// when
	err := _uut.Validate(request)

	// then
	assert.Nil(t, err)
}

func TestValidateUpdateUserRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.UpdateUserRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Name Test Cases
		{
			"Name is invalid because it's required",
			requests.UpdateUserRequest{Name: ""},
			"Name",
			"required",
		},
		{
			"Name is invalid because it has too many characters",
			requests.UpdateUserRequest{Name: strings.Repeat("a", 101)},
			"Name",
			"max",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateUserRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
