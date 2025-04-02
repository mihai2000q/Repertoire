package search

import (
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/internal"
)

type MeiliClient struct {
	meilisearch.ServiceManager
}

func NewMeiliClient(env internal.Env) MeiliClient {
	return MeiliClient{meilisearch.New(env.MeiliUrl, meilisearch.WithAPIKey(env.MeiliMasterKey))}
}
