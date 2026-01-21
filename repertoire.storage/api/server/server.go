package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"repertoire/storage/internal"

	"go.uber.org/fx"
)

func NewServer(lc fx.Lifecycle, handler *RequestHandler, env internal.Env) *http.Server {
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", env.ApplicationHost, env.ApplicationPort),
		Handler: handler.Gin,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return startServer(server)
		},
	})

	return server
}

func startServer(server *http.Server) error {
	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	fmt.Println("Starting HTTP server at", server.Addr)
	go func() {
		err := server.Serve(ln)
		if err != nil {
			log.Fatalf("Error starting the HTTP Server: %v", err)
		}
	}()
	return nil
}
