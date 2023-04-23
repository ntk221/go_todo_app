package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{
			Handler: mux,
		},
		l: l,
	}
}

func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// errgroup.WithContext()を使ってerrgroup.Groupを作成する
	eg, ctx := errgroup.WithContext(ctx)

	// TODO: errgroup.Go()を使って，別goroutineで，http.Server.ListenAndServe()を実行する
	eg.Go(func() error {
		if err := s.srv.Serve(s.l); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close server: %v", err)
			return err
		}
		return nil
	})

	// チャネルを使って，ctx.Done()が返す値を待つ
	<-ctx.Done()

	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}

	// errgroup.Group.Wait()を使って，errgroup.Groupが完了するのを待つ
	return eg.Wait()
}
