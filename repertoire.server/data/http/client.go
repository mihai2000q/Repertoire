package http

import (
	"github.com/go-resty/resty/v2"
	"repertoire/server/internal"
)

type Client interface {
	R() *resty.Request
}

func NewRestyClient(env internal.Env) Client {
	return resty.New().
		SetBaseURL(env.StorageUrl)
}
