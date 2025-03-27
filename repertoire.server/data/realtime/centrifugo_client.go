package realtime

import (
	"context"
	"github.com/centrifugal/centrifuge-go"
	"net/http"
	"repertoire/server/internal"
)

type CentrifugoClient interface {
	Connect() error
	Publish(ctx context.Context, channel string, data []byte) (centrifuge.PublishResult, error)
	Disconnect() error
	SetToken(token string)
}

func NewCentrifugoClient(env internal.Env) CentrifugoClient {
	return centrifuge.NewJsonClient(env.CentrifugoUrl, centrifuge.Config{
		Header: http.Header{
			"X-API-Key": []string{env.CentrifugoApiKey},
		},
	})
}
