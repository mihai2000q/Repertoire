package http

import (
	"github.com/go-resty/resty/v2"
	"repertoire/server/data/logger"
)

type Client struct {
	*resty.Client
}

func NewClient(logger *logger.RestyLogger) Client {
	return Client{
		resty.New().SetLogger(logger),
	}
}
