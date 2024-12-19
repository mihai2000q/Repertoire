package requests

import (
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/validation"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var validUserName = "Marcel Tembel"

func TestValidateUpdateUserRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := requests.UpdateUserRequest{
		Name: validUserName,
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
