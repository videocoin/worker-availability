package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/videocoin/worker-availablity/stats"
)

func report() {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		cancel()
	}()

	appctx, err := stats.NewContext(ctx, stats.FromEnv())
	if err != nil {
		appctx.Log.Fatalf("failed to bootstrap application %v", err)
	}
	if err := stats.Serve(appctx); err != nil {
		appctx.Log.Fatalf("serve crashed %v", err)
	}
}
