package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"

	"word-search-in-files/pkg/searcher"

	"word-search-in-files/pkg/server"
)

//go:embed examples/*
var dirWithFiles embed.FS

type configuration struct {
	Port            string        `envconfig:"PORT" default:"8080"`
	ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"5s"`
}

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {
	rootCtx := signalContext()

	var cfg configuration
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	s := &searcher.Searcher{FS: dirWithFiles}

	srv, err := server.NewServer(server.Config{
		Port:            cfg.Port,
		ShutdownTimeout: cfg.ShutdownTimeout,
		Searcher:        s,
	})
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	eg, ctx := errgroup.WithContext(rootCtx)

	eg.Go(func() error {
		return srv.Serve(ctx)
	})

	return eg.Wait()
}

// signalContext returns a context that is canceled if either SIGTERM or SIGINT signal is received.
func signalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-c
		slog.InfoContext(ctx, "received signal: %s", sig)
		cancel()
	}()

	return ctx
}
