package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"time"

	flag "github.com/spf13/pflag"
	"github.com/videocoin/worker-availablity/stats"
)

var (
	start *string = flag.String("start", "", "Start of the range in the RFC3339 format 92006-01-02T15:04:05Z07:00).")
	end   *string = flag.String("end", "", "End of the range in the RFC3339 format 92006-01-02T15:04:05Z07:00).")
	csv   *string = flag.String("csv", "", "Output file with csv report. (by default will be printed to stdout.")
)

func main() {
	flag.Parse()

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

	startT, err := time.Parse(time.RFC3339, *start)
	if err != nil {
		appctx.Log.Fatalf("failed to parse starting timestamp %v: %v", *start, err)
	}
	endT, err := time.Parse(time.RFC3339, *end)
	if err != nil {
		appctx.Log.Fatalf("failedd to parse ending timestamp %v: %v", *end, err)
	}

	var (
		f io.Writer
	)
	if len(*csv) > 0 {
		f, err = os.Create(*csv)
		if err != nil {
			appctx.Log.Fatalf("failed to create file at %v", *csv)
		}
	} else {
		f = os.Stdout
	}

	report, err := stats.CreateReport(appctx, startT, endT)
	if err != nil {
		appctx.Log.Fatalf("failed to create report %v", err)
	}

	if _, err := report.WriteTo(f); err != nil {
		appctx.Log.Fatalf("failed to save report to the file %v", err)
	}
	appctx.Log.Infof("application finished")
}
