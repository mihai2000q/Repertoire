package requests

import (
	"repertoire/server/api/validation"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var validUserName = "Marcel Tembel"

func TestValidateUpdateUserRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	// given
	_uut := validation.NewValidator(nil)

	request := UpdateUserRequest{
		ID:   uuid.New(),
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
		request              UpdateUserRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// ID Test Cases
		{
			"ID is invalid because it's required",
			UpdateUserRequest{ID: uuid.Nil, Name: validUserName},
			"ID",
			"required",
		},
		// Name Test Cases
		{
			"Name is invalid because it's required",
			UpdateUserRequest{ID: uuid.New(), Name: ""},
			"Name",
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
			assert.Contains(t, errCode.Error.Error(), "UpdateUserRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, 400, errCode.Code)
		})
	}
}
