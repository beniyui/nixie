package general

import (
	"context"
	"os/signal"
	"syscall"
)

func stay() {

	//プログラムの終了を検知するクライアント
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//ここでストップする
	select {
	case <-ctx.Done():
	}
}
