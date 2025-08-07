package message

import (
	"context"
	"repertoire/server/data/logger"
	"repertoire/server/data/service"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/domain/message/handler/playlist"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/domain/message/handler/song"
	"repertoire/server/domain/message/handler/storage"
	"repertoire/server/domain/message/handler/user"
	"repertoire/server/internal/message/topics"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/cenkalti/backoff/v4"

	"log"
	"repertoire/server/domain/message/handler/artist"
	"time"

	"go.uber.org/fx"
)

type messageHandler interface {
	Handle(msg *message.Message) error
	GetName() string
	GetTopic() topics.Topic
}

func NewRouter(
	lc fx.Lifecycle,
	messagePublisherService service.MessagePublisherService,
	logger *logger.WatermillLogger,

	albumCreatedHandler album.AlbumCreatedHandler,
	albumDeletedHandler album.AlbumsDeletedHandler,
	albumsUpdatedHandler album.AlbumsUpdatedHandler,

	artistCreatedHandler artist.ArtistCreatedHandler,
	artistDeletedHandler artist.ArtistsDeletedHandler,
	artistUpdatedHandler artist.ArtistUpdatedHandler,

	playlistCreatedHandler playlist.PlaylistCreatedHandler,
	playlistDeletedHandler playlist.PlaylistDeletedHandler,
	playlistUpdatedHandler playlist.PlaylistUpdatedHandler,

	songCreatedHandler song.SongCreatedHandler,
	songDeletedHandler song.SongsDeletedHandler,
	songsUpdatedHandler song.SongsUpdatedHandler,

	userDeletedHandler user.UserDeletedHandler,

	addToSearchEngineHandler search.AddToSearchEngineHandler,
	deleteFromSearchEngineHandler search.DeleteFromSearchEngineHandler,
	updateFromSearchEngineHandler search.UpdateFromSearchEngineHandler,

	deleteDirectoriesStorageHandler storage.DeleteDirectoriesStorageHandler,
) *message.Router {
	handlers := []messageHandler{
		albumCreatedHandler,
		albumDeletedHandler,
		albumsUpdatedHandler,

		artistCreatedHandler,
		artistDeletedHandler,
		artistUpdatedHandler,

		playlistCreatedHandler,
		playlistDeletedHandler,
		playlistUpdatedHandler,

		songCreatedHandler,
		songDeletedHandler,
		songsUpdatedHandler,

		userDeletedHandler,

		addToSearchEngineHandler,
		deleteFromSearchEngineHandler,
		updateFromSearchEngineHandler,

		deleteDirectoriesStorageHandler,
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		log.Fatal(err)
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		CustomRetryMiddleware{
			MaxRetries:      2,
			InitialInterval: time.Millisecond * 100,
			Logger:          logger,
		}.Middleware,
		middleware.Recoverer,
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			for _, handler := range handlers {
				router.AddNoPublisherHandler(
					handler.GetName(),
					string(topics.TopicToQueueMap[handler.GetTopic()]),
					messagePublisherService.GetClient(),
					func(msg *message.Message) error {
						topic := msg.Metadata.Get("topic")
						if topic != string(handler.GetTopic()) {
							return nil
						}
						return handler.Handle(msg)
					},
				)
			}

			go func() {
				if err := router.Run(context.Background()); err != nil {
					// Log the error and stop the FX app
					log.Fatalf("Router stopped with error: %v", err)
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			return router.Close()
		},
	})

	return router
}

// Copied the Retry Middleware implementation and ACKNOWLEDGED the message on the last retry
// to fix infinite loop of retrying the message

type CustomRetryMiddleware struct {
	MaxRetries      int
	InitialInterval time.Duration
	Logger          watermill.LoggerAdapter
}

func (c CustomRetryMiddleware) Middleware(h message.HandlerFunc) message.HandlerFunc {
	return func(msg *message.Message) ([]*message.Message, error) {
		producedMessages, err := h(msg)
		if err == nil {
			return producedMessages, nil
		}

		expBackoff := backoff.NewExponentialBackOff()
		expBackoff.InitialInterval = c.InitialInterval

		ctx := msg.Context()

		retryNum := 1
		expBackoff.Reset()
	retryLoop:
		for {
			waitTime := expBackoff.NextBackOff()
			select {
			case <-ctx.Done():
				msg.Ack() // Acknowledge the message to stop retrying
				return producedMessages, err
			case <-time.After(waitTime):
				// go on
			}

			producedMessages, err = h(msg)
			if err == nil {
				return producedMessages, nil
			}

			if c.Logger != nil {
				c.Logger.Error("Error occurred, retrying", err, watermill.LogFields{
					"retry_no":     retryNum,
					"max_retries":  c.MaxRetries,
					"wait_time":    waitTime,
					"elapsed_time": expBackoff.GetElapsedTime(),
				})
			}

			retryNum++
			if retryNum > c.MaxRetries {
				msg.Ack() // Acknowledge the message to stop retrying
				break retryLoop
			}
		}

		return nil, err
	}
}
