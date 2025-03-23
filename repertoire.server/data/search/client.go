package search

import (
	"github.com/meilisearch/meilisearch-go"
	"repertoire/server/internal"
)

type Client struct {
	meilisearch.ServiceManager
}

func NewMeiliClient(env internal.Env) Client {
	url := "http://" + env.MeiliHost + ":" + env.MeiliPort
	return Client{meilisearch.New(url, meilisearch.WithAPIKey(env.MeiliMasterKey))}
}
