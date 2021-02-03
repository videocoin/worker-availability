package stats

import (
	"os"
	"bufio"
	"context"
	"fmt"
	"io"
	"time"
)

func (r Report) ReportIncentives(w io.Writer) (total int, err error) {
	var (
		n int
	)
	const incentivePerHour = 0.14 / (365.0 * 24.0)

	n, err = fmt.Fprintln(w, "worker_address,incentive")
	total += n
	if err != nil {
		return
	}
	for _, info := range r {
		itemCount := len(info.Configuration)
		if(itemCount > 0 && len(info.Address) > 0) {
			incentive := 0.0 
			for _, conf := range info.Configuration {
				incentive += conf.Online.Hours() * conf.DirectStake * incentivePerHour
			}
			n, err = fmt.Fprintf(w, "%s,%v\n", info.Address, incentive)
			total += n
		}
	}
	return
}

func CreateIncentives(appctx Context, ctx context.Context, fileIncentives string, fileUptime string, start time.Time, end time.Time) error {
	rep, err := CreateReport(appctx, ctx, start, end)
	if err != nil {
        return err
    }
	f, err := os.Create(fileIncentives)
	if err != nil {
        return err
    }
	w := bufio.NewWriter(f)
	rep.ReportIncentives(w)
	w.Flush()
	// Create report
	f2, err := os.Create(fileUptime)
	if err != nil {
        return err
    }
	w2 := bufio.NewWriter(f2)
	rep.WriteTo(w2)
	w2.Flush()

	return nil
}
