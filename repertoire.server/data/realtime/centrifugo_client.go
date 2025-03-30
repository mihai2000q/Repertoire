package realtime

import (
	"context"
	"github.com/centrifugal/centrifuge-go"
	"repertoire/server/internal"
)

type CentrifugoClient interface {
	Connect() error
	Publish(ctx context.Context, channel string, data []byte) (centrifuge.PublishResult, error)
	Disconnect() error
	SetToken(token string)
}

func NewCentrifugoClient(env internal.Env) CentrifugoClient {
	return centrifuge.NewJsonClient(env.CentrifugoUrl, centrifuge.Config{})
}
