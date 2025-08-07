package search

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"repertoire/server/internal/enums"
	"repertoire/server/internal/wrapper"
	"repertoire/server/model"
	"repertoire/server/test/integration/test/core"
	searchData "repertoire/server/test/integration/test/data/search"
	"repertoire/server/test/integration/test/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchGet_WhenSuccessful_ShouldReturnSearchResults(t *testing.T) {
	// given
	utils.SeedAndCleanupSearchData(t, searchData.GetSearchDocuments())

	query := "justice"

	expectedResults := []any{
		searchData.SongSearches[0],
		searchData.AlbumSearches[0],
		searchData.SongSearches[1],
		searchData.PlaylistSearches[1],
	}
	expectedTotalCount := len(expectedResults)

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithUser(model.User{ID: searchData.UserID}).
		GET(w, "/api/search?query="+query)

	// then
	assert.Equal(t, http.StatusOK, w.Code)

	var response wrapper.WithTotalCount[any]
	_ = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, int64(expectedTotalCount), response.TotalCount)
	assert.Len(t, response.Models, len(expectedResults))
	for i, expectedResult := range expectedResults {
		curr := response.Models[i]
		expectedID := utils.UnmarshallDocument[map[string]any](expectedResult)["id"]
		actualID := curr.(map[string]interface{})["id"]
		var prefix string
		switch curr.(map[string]interface{})["type"] {
		case string(enums.Artist):
			prefix = "artist-"
		case string(enums.Album):
			prefix = "album-"
		case string(enums.Song):
			prefix = "song-"
		case string(enums.Playlist):
			prefix = "playlist-"
		}

		assert.Equal(t, expectedID, prefix+actualID.(string))
	}
}
