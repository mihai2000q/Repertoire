package search

import (
	"repertoire/server/internal/env"

	"github.com/meilisearch/meilisearch-go"
)

type MeiliClient struct {
	meilisearch.ServiceManager
}

func NewMeiliClient(env env.Env) MeiliClient {
	return MeiliClient{meilisearch.New(env.MeiliUrl, meilisearch.WithAPIKey(env.MeiliMasterKey))}
}
