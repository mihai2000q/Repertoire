package internal

import "github.com/go-resty/resty/v2"

func NewRestyClient(env Env) *resty.Client {
	return resty.New().
		SetBaseURL(env.StorageUrl)
}
