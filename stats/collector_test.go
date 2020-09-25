package stats

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
)

func TestCollector(t *testing.T) {
	appctx, err := NewContext(context.TODO(), FromEnv())
	require.NoError(t, err)
	Collect(appctx)

	records, err := appctx.DB.All(context.TODO())
	require.NoError(t, err)
	require.True(t, len(records) > 0)
	for _, record := range records {
		require.Equal(t, emitterv1.WorkerStateBonded, record.Miner.WorkerState)
	}
}
