package message

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"repertoire/server/data/service"
	"repertoire/server/domain/message/handler/album"
	"repertoire/server/domain/message/handler/playlist"
	"repertoire/server/domain/message/handler/search"
	"repertoire/server/domain/message/handler/song"
	"repertoire/server/domain/message/handler/user"
	"repertoire/server/internal/message/topics"

	"go.uber.org/fx"
	"log"
	"repertoire/server/domain/message/handler/artist"
	"time"
)

type messageHandler interface {
	Handle(msg *message.Message) error
	GetName() string
	GetTopic() topics.Topic
}

func NewRouter(
	lc fx.Lifecycle,
	messagePublisherService service.MessagePublisherService,

	albumCreatedHandler album.AlbumCreatedHandler,
	albumDeletedHandler album.AlbumDeletedHandler,
	albumsUpdatedHandler album.AlbumsUpdatedHandler,

	artistCreatedHandler artist.ArtistCreatedHandler,
	artistDeletedHandler artist.ArtistDeletedHandler,
	artistUpdatedHandler artist.ArtistUpdatedHandler,

	playlistCreatedHandler playlist.PlaylistCreatedHandler,
	playlistDeletedHandler playlist.PlaylistDeletedHandler,
	playlistUpdatedHandler playlist.PlaylistUpdatedHandler,

	songCreatedHandler song.SongCreatedHandler,
	songDeletedHandler song.SongDeletedHandler,
	songsUpdatedHandler song.SongsUpdatedHandler,

	userDeletedHandler user.UserDeletedHandler,

	addToSearchEngineHandler search.AddToSearchEngineHandler,
	deleteFromSearchEngineHandler search.DeleteFromSearchEngineHandler,
	updateFromSearchEngineHandler search.UpdateFromSearchEngineHandler,
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
	}

	logger := watermill.NewStdLogger(false, false)
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		log.Fatal(err)
	}

	router.AddMiddleware(
		middleware.Retry{
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
