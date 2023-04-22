package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	s.ListenAndServe()
	// TODO: context.Contextを通じて処理のキャンセルを受け取った時，
	// *http.Server.Shutdown()を呼び出してサーバーを終了させる

	// TODO: run 関数の戻り値をerror型に変更し，エラーを返すようにする
}

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failet to terminate server: %v", err)
	}
}
