package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type searcher interface {
	Search(word string) (files []string, err error)
}

type Config struct {
	Port            string        `validate:"required"`
	Searcher        searcher      `validate:"required"`
	ShutdownTimeout time.Duration `validate:"required"`
}

type Server struct {
	srv      *http.Server
	cfg      Config
	searcher searcher
}

func NewServer(cfg Config) (*Server, error) {
	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, fmt.Errorf("config validation error: %w", err)
	}

	srv := &http.Server{
		Addr: ":" + cfg.Port,
	}

	s := &Server{
		srv:      srv,
		cfg:      cfg,
		searcher: cfg.Searcher,
	}
	srv.Handler = s.initServerHandler()

	return s, nil
}

func (s *Server) initServerHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/search", s.handleSearch)
	return mux
}

// Serve starts HTTP server and stops it when the provided context is canceled.
func (s *Server) Serve(ctx context.Context) error {
	errChan := make(chan error, 1)
	defer close(errChan)
	go func() {
		errChan <- s.srv.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err

	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := s.srv.Shutdown(ctxShutdown); err != nil { //nolint:contextcheck // shutdown had to be executed with independent context
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return nil
	}
}
