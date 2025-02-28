package search

import (
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/internal"
)

func NewMeiliClient(env internal.Env) meilisearch.ServiceManager {
	return meilisearch.New(env.MeiliUrl, meilisearch.WithAPIKey(env.MeiliMasterKey))
}
