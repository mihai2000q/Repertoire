package http

import (
	"repertoire/server/data/logger"
	"repertoire/server/internal"

	"github.com/go-resty/resty/v2"
)

type RestyClient struct {
	*resty.Client
}

func NewRestyClient(logger *logger.RestyLogger, env internal.Env) RestyClient {
	return RestyClient{
		resty.New().
			SetLogger(logger).
			SetDebugBodyLimit(1024).
			SetDebug(env.LogLevel == internal.DebugLogLevel),
	}
}
