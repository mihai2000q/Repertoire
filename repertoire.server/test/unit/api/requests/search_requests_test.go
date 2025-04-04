package requests

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"repertoire/server/api/requests"
	"repertoire/server/api/validation"
	"repertoire/server/internal/enums"
	"testing"
)

func TestValidateSearchGetRequest_WhenIsValid_ShouldReturnNil(t *testing.T) {
	tests := []struct {
		name    string
		request requests.SearchGetRequest
	}{
		{
			"Empty",
			requests.SearchGetRequest{},
		},
		{
			"With Pagination",
			requests.SearchGetRequest{
				CurrentPage: &[]int{1}[0],
				PageSize:    &[]int{20}[0],
			},
		},
		{
			"With Type",
			requests.SearchGetRequest{
				Type: &[]enums.SearchType{enums.Song}[0],
			},
		},
		{
			"With Filtering and Sorting",
			requests.SearchGetRequest{
				Filter: []string{"release_date > 145023"},
				Order:  []string{"release_date asc"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// given
			_uut := validation.NewValidator(nil)

			// when
			errCode := _uut.Validate(tt.request)

			// then
			assert.Nil(t, errCode)
		})
	}
}

func TestValidateSearchGetRequest_WhenSingleFieldIsInvalid_ShouldReturnBadRequest(t *testing.T) {
	tests := []struct {
		name                 string
		request              requests.SearchGetRequest
		expectedInvalidField string
		expectedFailedTag    string
	}{
		// Current Page
		{
			"Current Page is invalid because page size is required with Page Size",
			requests.SearchGetRequest{
				PageSize: &[]int{20}[0],
			},
			"CurrentPage",
			"required_with",
		},
		// Page Size
		{
			"Page Size is invalid because page size is required with Current Page",
			requests.SearchGetRequest{
				CurrentPage: &[]int{1}[0],
			},
			"PageSize",
			"required_with",
		},
		// Type
		{
			"Type is invalid because it is not part of the enum",
			requests.SearchGetRequest{
				Type: &[]enums.SearchType{"something"}[0],
			},
			"Type",
			"search_type_enum",
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
			assert.Contains(t, errCode.Error.Error(), "SearchGetRequest."+tt.expectedInvalidField)
			assert.Contains(t, errCode.Error.Error(), "'"+tt.expectedFailedTag+"' tag")
			assert.Equal(t, http.StatusBadRequest, errCode.Code)
		})
	}
}
