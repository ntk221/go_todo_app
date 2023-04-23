package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	t.Skip("リファクタリング中")
	// TEST用のの環境を用意する
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})

	// テスト用のHTTPリクエストを送る
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Fatal(err)
	}

	// responseがある場合は，deferでBodyを閉じる
	defer rsp.Body.Close()
	// responseのBodyを読み込む
	got, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read Body:%v", err)
	}
	// responseのBodyの内容をテストする
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %s, but %s:", want, got)
	}

	// テストが終わったら，contextをキャンセルする
	cancel()

	// run関数の戻り値を検証
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}
