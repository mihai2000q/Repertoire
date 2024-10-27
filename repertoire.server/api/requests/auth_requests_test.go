package requests

import (
	"github.com/stretchr/testify/assert"
	"repertoire/api/validation"
	"strings"
	"testing"
)

var validEmail = "someone@yahoo.com"
var validPassword = "Password123"

func TestValidateSignInRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := SignInRequest{
		Email:    validEmail,
		Password: validPassword,
	}

	errCode := _uut.Validate(request)

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
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			err := errCode.Error.Error()

			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "SignInRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}

var validName = "Samuel"

func TestValidateSignUpRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	_uut := validation.NewValidator(nil)

	request := SignUpRequest{
		Name:     validName,
		Email:    validEmail,
		Password: validPassword,
	}

	errCode := _uut.Validate(request)

	assert.Nil(t, errCode)
}

func TestValidateSignUpRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              SignUpRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Name Test Cases
		{
			"Name is invalid because it's required",
			SignUpRequest{Name: "", Email: validEmail, Password: validPassword},
			"Name",
			"required",
		},
		// Email Test Cases
		{
			"Email is invalid because it's required",
			SignUpRequest{Name: validName, Email: "", Password: validPassword},
			"Email",
			"required",
		},
		{
			"Email is invalid because it has too many characters",
			SignUpRequest{Name: validName, Email: strings.Repeat("a", 257), Password: validPassword},
			"Email",
			"max",
		},
		{
			"Email is invalid because it is not an email",
			SignUpRequest{Name: validName, Email: "someone", Password: validPassword},
			"Email",
			"email",
		},
		{
			"Email is invalid because it is not an email",
			SignUpRequest{Name: validName, Email: "someone@yahoo", Password: validPassword},
			"Email",
			"email",
		},
		{
			"Email is invalid because it is not an email",
			SignUpRequest{Name: validName, Email: "someone.com", Password: validPassword},
			"Email",
			"email",
		},
		// Password Test Cases
		{
			"Password is invalid because it's required",
			SignUpRequest{Name: validName, Email: validEmail, Password: ""},
			"Password",
			"required",
		},
		{
			"Password is invalid because it has less than 8 characters",
			SignUpRequest{Name: validName, Email: validEmail, Password: "1234567"},
			"Password",
			"min",
		},
		{
			"Password is invalid because it doesn't have an uppercase letter",
			SignUpRequest{Name: validName, Email: validEmail, Password: strings.Repeat("a", 9)},
			"Password",
			"hasUpper",
		},
		{
			"Password is invalid because it doesn't have a lowercase letter",
			SignUpRequest{Name: validName, Email: validEmail, Password: strings.Repeat("A", 9)},
			"Password",
			"hasLower",
		},
		{
			"Password is invalid because it doesn't have any digit",
			SignUpRequest{Name: validName, Email: validEmail, Password: strings.Repeat("A", 4) + strings.Repeat("a", 4)},
			"Password",
			"hasDigit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_uut := validation.NewValidator(nil)

			errCode := _uut.Validate(tt.request)

			err := errCode.Error.Error()

			assert.NotNil(t, errCode)
			assert.Len(t, errCode.Error, 1)
			assert.Contains(t, err, "SignUpRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
