package stats

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

type Report map[string]*WorkerInfo

func (r Report) WriteTo(w io.Writer) (total int, err error) {
	var (
		n int
	)
	n, err = fmt.Fprintln(w, "worker,client_id,worker_address,configuration_hash,duration_online")
	if err != nil {
		return
	}
	total += n
	for _, info := range r {
		for _, conf := range info.Configuration {
			n, err = fmt.Fprintf(w, "%s,%s,%s,0x%x,%v\n", info.Name, info.ClientID, info.Address, conf.Hash, conf.Online)
			if err != nil {
				return
			}
			total += n
		}
	}
	return
}

type WorkerInfo struct {
	Name          string
	ClientID      string
	Address       string
	Configuration []*ConfigurationInfo
}

type ConfigurationInfo struct {
	Hash   []byte
	Online time.Duration
}

func CreateReport(appctx Context, start, end time.Time) (Report, error) {
	var (
		rep       = Report{}
		err       error
		timestamp time.Time
		init      = false
		hasher    = sha256.New()
	)
	err = appctx.DB.Process(appctx, start, end, func(record Aggregated) bool {
		if !init {
			timestamp = record.Timestamp
			init = true
			return true
		}
		d := record.Timestamp.Sub(timestamp)
		appctx.Log.Debugf("delta between %v and %v - %v", record.Timestamp, timestamp, d)
		timestamp = record.Timestamp
		for _, miner := range record.Records {
			err = binary.Write(hasher, binary.LittleEndian, miner.Miner.SystemInfo.CpuCores)
			if err != nil {
				return false
			}
			err = binary.Write(hasher, binary.LittleEndian, miner.Miner.SystemInfo.CpuFreq)
			if err != nil {
				return false
			}
			err = binary.Write(hasher, binary.LittleEndian, miner.Miner.SystemInfo.MemTotal)
			if err != nil {
				return false
			}

			hash := make([]byte, 0, 32)
			hash = hasher.Sum(hash)
			hasher.Reset()

			info, exist := rep[miner.Miner.Name]
			if !exist {
				info = &WorkerInfo{}
				rep[miner.Miner.Name] = info
				info.Name = miner.Miner.Name
				info.ClientID = miner.Miner.Id
				info.Address = miner.Miner.Address
				info.Configuration = []*ConfigurationInfo{
					{
						Hash:   hash,
						Online: d,
					},
				}
			} else {
				last := len(info.Configuration) - 1
				if bytes.Compare(info.Configuration[last].Hash, hash) == 0 {
					info.Configuration[last].Online += d
				} else {
					appctx.Log.Debugf("observed change of the configuration for worker %v", miner.Miner.Name)
					info.Configuration = append(info.Configuration,
						&ConfigurationInfo{Hash: hash, Online: d})
				}
			}

		}
		return true
	})
	return rep, err
}
