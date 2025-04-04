package server

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"log"
	"net"
	"net/http"
	"os"
	"repertoire/server/data/logger"
	"repertoire/server/internal"
)

func NewServer(lc fx.Lifecycle, handler *RequestHandler, logger *logger.Logger, env internal.Env) *http.Server {
	address := ""
	if os.Getenv("INTEGRATION_TESTING_ENVIRONMENT_FILE_PATH") == "" {
		address = fmt.Sprintf("%s:%s", env.ApplicationHost, env.ApplicationPort)
	}

	server := &http.Server{
		Addr:    address,
		Handler: handler.Gin,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return startServer(server, logger)
		},
	})

	return server
}

func startServer(server *http.Server, logger *logger.Logger) error {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	logger.Info("Starting HTTP server at " + server.Addr)
	go func() {
		err = server.Serve(ln)
		if err != nil {
			log.Fatalf("Error starting the HTTP Server: %v", err)
		}
	}()
	return nil
}
