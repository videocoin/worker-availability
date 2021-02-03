package main

import (
	"os"
	"os/signal"
	"context"
	"time"

	"github.com/videocoin/worker-availablity/stats"
)

func incentives(fileIncentives string, fileUptime string, startTime  string, endTime  string) {
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

	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		appctx.Log.Fatalf("failed to parse start time %v", err)
	}
	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		appctx.Log.Fatalf("failed to parse end time %v", err)
	}

	if err := stats.CreateIncentives(appctx, ctx, fileIncentives, fileUptime, start, end); err != nil {
		appctx.Log.Fatalf("incentive creation failed %v", err)
	}	
}
