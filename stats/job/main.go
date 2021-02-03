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
	} else if len(os.Args) > 5 && os.Args[1] == "incentives" {
		fileIncentives := os.Args[2]
		fileUptime := os.Args[3]
		startTime := os.Args[4]
		endTime := os.Args[5]
		incentives(fileIncentives, fileUptime, startTime, endTime)
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
