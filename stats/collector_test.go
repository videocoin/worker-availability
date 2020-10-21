package stats

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	v1 "github.com/videocoin/cloud-api/miners/v1"
)

func TestCollector(t *testing.T) {
	appctx, err := NewContext(context.TODO(), FromEnv())
	require.NoError(t, err)
	Collect(appctx)
	Collect(appctx)

	records, err := appctx.DB.All(context.TODO())
	require.NoError(t, err)
	require.True(t, len(records) > 0)
	for _, record := range records {
		require.Equal(t, emitterv1.WorkerStateBonded, record.Miner.WorkerState)
		fmt.Println(record.Miner.Name, record.Timestamp)
	}
}

func TestReport(t *testing.T) {
	now := time.Now()
	t1 := now.Add(time.Hour)
	t2 := t1.Add(time.Hour)
	records := []interface{}{
		Record{Timestamp: now, Miner: &v1.MinerResponse{Name: "first", Id: "1", SystemInfo: &v1.SystemInfo{}}},
		Record{Timestamp: t1, Miner: &v1.MinerResponse{Name: "first", Id: "1", SystemInfo: &v1.SystemInfo{}}},
		Record{Timestamp: t1, Miner: &v1.MinerResponse{Name: "second", Id: "2", SystemInfo: &v1.SystemInfo{}}},
		Record{Timestamp: t2, Miner: &v1.MinerResponse{Name: "first", Id: "1", SystemInfo: &v1.SystemInfo{}}},
		Record{Timestamp: t2, Miner: &v1.MinerResponse{Name: "second", Id: "2", SystemInfo: &v1.SystemInfo{}}},
	}
	appctx, err := NewContext(context.TODO(), FromEnv())
	require.NoError(t, err)

	require.NoError(t, appctx.DB.Save(context.TODO(), records))
	report, err := CreateReport(appctx, now, t2)
	require.NoError(t, err)
	report.WriteTo(os.Stdout)
}
