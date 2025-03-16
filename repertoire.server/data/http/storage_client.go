package http

import (
	"github.com/go-resty/resty/v2"
	"repertoire/server/internal"
)

type StorageClient interface {
	R() *resty.Request
}

func NewRestyClient(env internal.Env) StorageClient {
	return resty.New().SetBaseURL(env.StorageUrl)
}
