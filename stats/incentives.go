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
	n, err = fmt.Fprintln(w, "worker_address,incentive")
	if err != nil {
		return
	}
	total += n
	for _, info := range r {
		itemCount := len(info.Configuration)
		if(itemCount > 0 && len(info.Address) > 0) {
			conf := info.Configuration[itemCount - 1]
			n, err = fmt.Fprintf(w, "%s,%v\n", info.Address, conf.Incentive)
			if err != nil {
				return
			}
			total += n
		}
	}
	return
}

func CreateIncentives(appctx Context, ctx context.Context, fileName string, start time.Time, end time.Time) error {
	rep, err := CreateReport(appctx, ctx, start, end)
	if err != nil {
        return err
    }
	f, err := os.Create(fileName)
	if err != nil {
        return err
    }
	w := bufio.NewWriter(f)
	rep.ReportIncentives(w)
	w.Flush()
	return nil
}
