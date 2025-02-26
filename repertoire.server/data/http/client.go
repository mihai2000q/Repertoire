package http

import (
	"github.com/go-resty/resty/v2"
	"repertoire/server/internal"
)

func NewRestyClient(env internal.Env) *resty.Client {
	return resty.New().
		SetBaseURL(env.StorageUrl)
}
