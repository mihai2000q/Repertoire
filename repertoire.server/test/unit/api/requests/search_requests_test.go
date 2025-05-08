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
			"With Filtering",
			requests.SearchGetRequest{
				Filter: []string{
					"release_date > 145023",
					"(album IS NULL OR title = \"something else\" OR artist.id = 12345)",
					"NOT release_date > 123",
					"NOT release_date IS NULL",
					"release_date IS NOT NULL",
					"(artist IS NULL)",
					"id IN [12, 123]",
				},
			},
		},
		{
			"With Sorting",
			requests.SearchGetRequest{
				Order: []string{"release_date:asc", "price:desc"},
			},
		},
		{
			"With IDs",
			requests.SearchGetRequest{
				Type: &[]enums.SearchType{enums.Song}[0],
				IDs:  []string{"1", "2"},
			},
		},
		{
			"With Not IDs",
			requests.SearchGetRequest{
				Type:   &[]enums.SearchType{enums.Song}[0],
				NotIDs: []string{"1", "2"},
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
			"Current Page is invalid because it is required with Page Size",
			requests.SearchGetRequest{
				PageSize: &[]int{20}[0],
			},
			"CurrentPage",
			"required_with",
		},
		// Page Size
		{
			"Page Size is invalid because it is required with Current Page",
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
		// IDs
		{
			"IDs is invalid because the type is required with",
			requests.SearchGetRequest{
				IDs: []string{"1"},
			},
			"Type",
			"required_with",
		},
		// Not IDs
		{
			"NotIDs is invalid because the type is required with",
			requests.SearchGetRequest{
				NotIDs: []string{"1"},
			},
			"Type",
			"required_with",
		},
		// Filter
		{
			"Filter is invalid",
			requests.SearchGetRequest{
				Filter: []string{"price something else"},
			},
			"Filter",
			"search_filter",
		},
		{
			"Filter is invalid",
			requests.SearchGetRequest{
				Filter: []string{"price = else or price = some"},
			},
			"Filter",
			"search_filter",
		},
		// Order
		{
			"Order is invalid because of the separator",
			requests.SearchGetRequest{
				Order: []string{"price asc"},
			},
			"Order",
			"search_order",
		},
		{
			"Order is invalid because of the order type",
			requests.SearchGetRequest{
				Order: []string{"price:ascending"},
			},
			"Order",
			"search_order",
		},
		{
			"Order is invalid because of the spaces",
			requests.SearchGetRequest{
				Order: []string{"the price is   :ascending"},
			},
			"Order",
			"search_order",
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
