package http

import (
	"github.com/go-resty/resty/v2"
	"repertoire/server/internal"
)

type StorageClient struct {
	*resty.Client
}

func NewStorageClient(client Client, env internal.Env) StorageClient {
	return StorageClient{
		client.SetBaseURL(env.StorageUrl),
	}
}
