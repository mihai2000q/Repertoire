package search

import (
	"repertoire/server/internal"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliClient struct {
	meilisearch.ServiceManager
}

func NewMeiliClient(env internal.Env) MeiliClient {
	return MeiliClient{meilisearch.New(env.MeiliUrl, meilisearch.WithAPIKey(env.MeiliMasterKey))}
}
