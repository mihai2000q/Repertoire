package http

import (
	"github.com/go-resty/resty/v2"
	"repertoire/server/data/logger"
	"repertoire/server/internal"
)

type Client struct {
	*resty.Client
}

func NewClient(logger *logger.RestyLogger, env internal.Env) Client {
	return Client{
		resty.New().
			SetLogger(logger).
			SetDebug(env.LogLevel == internal.DebugLogLevel),
	}
}
