package search

import (
	"encoding/json"
	"github.com/centrifugal/centrifuge-go"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"repertoire/server/test/integration/test/core"
	"repertoire/server/test/integration/test/utils"
	"strconv"
	"testing"
)

func TestSearchMeiliWebhook_WhenSuccessful_ShouldPublishOnCentrifugo(t *testing.T) {
	// given
	// prepare cache
	taskID := 23
	userID := "some-user-id"
	core.MeiliCache.Set("task-"+strconv.Itoa(taskID), userID, cache.DefaultExpiration)

	// prepare centrifugo client assertion
	centrifugoClient := utils.GetCentrifugoClient(t)
	sub, _ := centrifugoClient.NewSubscription("search:" + userID)
	_ = centrifugoClient.Connect()
	_ = sub.Subscribe()
	sub.OnPublication(func(event centrifuge.PublicationEvent) {
		var data map[string]interface{}
		_ = json.Unmarshal(event.Data, &data)
		assert.Equal(t, "SEARCH_CACHE_INVALIDATION", data["action"])
		_ = sub.Unsubscribe()
		_ = centrifugoClient.Disconnect()
	})

	var task = struct {
		UID    int64  `json:"uid"`
		Status string `json:"status"`
	}{
		UID:    int64(taskID),
		Status: "succeeded",
	}

	// when
	w := httptest.NewRecorder()
	core.NewTestHandler().
		WithMeiliAuthentication().
		POSTZipped(w, "/api/search/meili-webhook", task)

	// then
	assert.Equal(t, http.StatusOK, w.Code)
}
