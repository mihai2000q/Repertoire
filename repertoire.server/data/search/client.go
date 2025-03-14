package search

import (
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/internal"
)

func NewMeiliClient(env internal.Env) meilisearch.ServiceManager {
	url := "http://" + env.MeiliHost + ":" + env.MeiliPort
	return meilisearch.New(url, meilisearch.WithAPIKey(env.MeiliMasterKey))
}
