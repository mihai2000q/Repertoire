package realtime

import (
	"context"
	"repertoire/server/internal"

	"github.com/centrifugal/centrifuge-go"
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
