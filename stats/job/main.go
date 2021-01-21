package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/videocoin/worker-availablity/stats"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "report" {
		report()
		return
	} else if len(os.Args) > 4 && os.Args[1] == "incentives" {
		fileName := os.Args[2]
		startTime := os.Args[3]
		endTime := os.Args[4]
		incentives(fileName, startTime, endTime)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		cancel()
	}()
	appctx, err := stats.NewContext(ctx, stats.FromEnv())
	if err != nil {
		appctx.Log.Errorf("failed to bootstrap application %v", err)
		os.Exit(1)
	}
	stats.Poll(appctx, stats.Collect)
	appctx.Log.Infof("application was stopped")
}
