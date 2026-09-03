package http

import (
	"repertoire/server/data/logger"
	"repertoire/server/internal/env"

	"github.com/go-resty/resty/v2"
)

type RestyClient struct {
	*resty.Client
}

func NewRestyClient(logger *logger.RestyLogger, env env.Env) RestyClient {
	return RestyClient{
		resty.New().
			SetLogger(logger).
			SetDebugBodyLimit(1024).
			SetDebug(env.LogLevel == env.DebugLogLevel),
	}
}
