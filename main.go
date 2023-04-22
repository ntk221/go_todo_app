package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	// errgroup.WithContext()を使ってerrgroup.Groupを作成する
	eg, ctx := errgroup.WithContext(ctx)

	// TODO: errgroup.Go()を使って，別goroutineで，http.Server.ListenAndServe()を実行する
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close server: %v", err)
			return err
		}
		return nil
	})

	// チャネルを使って，ctx.Done()が返す値を待つ
	<-ctx.Done()
	if err := s.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown server: %v", err)
	}

	// errgroup.Group.Wait()を使って，errgroup.Groupが完了するのを待つ
	return eg.Wait()
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failet to terminate server: %v", err)
	}
}
