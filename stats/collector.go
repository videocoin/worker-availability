package stats

import (
	"context"
	"fmt"
	"net/http"
	"time"

	stackdriver "github.com/TV4/logrus-stackdriver-formatter"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	emitterv1 "github.com/videocoin/cloud-api/emitter/v1"
	v1 "github.com/videocoin/cloud-api/miners/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Service                = "availability-job"
	AvailabilityCollection = "availability"
	EnvNamespace           = "AJ"
)

func FromEnv() (cfg Config) {
	if err := envconfig.Process(EnvNamespace, &cfg); err != nil {
		panic("failed to process enconfig " + err.Error())
	}
	return cfg
}

type Config struct {
	URL      string        `default:"https://console.videocoin.network/api/v1/miners/all"`
	Mongo    string        `default:"mongodb://localhost:27017"`
	Database string        `default:"availability"`
	Period   time.Duration `default:"1m"`
	Timeout  time.Duration `default:"1m"`
	Retries  int           `default:"3"`
	LogLevel string        `default:"debug"`
}

func NewContext(ctx context.Context, cfg Config) (appctx Context, err error) {
	logger := logrus.New()
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return appctx, fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}
	logger.SetLevel(logLevel)
	logger.SetFormatter(stackdriver.NewFormatter(
		stackdriver.WithService(Service),
	))
	log := logrus.NewEntry(logger)

	client := &http.Client{Timeout: cfg.Timeout}
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo))
	if err != nil {
		return appctx, err
	}
	return Context{
		Context: ctx,
		C:       cfg,
		Log:     log,
		HTTP:    client,
		DB: DB{Client: conn.Database(cfg.Database).Collection(
			AvailabilityCollection, options.Collection())},
	}, nil
}

type Context struct {
	context.Context
	C    Config
	Log  *logrus.Entry
	HTTP *http.Client
	DB   DB
}

func Poll(appctx Context, f func(Context)) {
	var (
		ticker = time.NewTicker(appctx.C.Period)
		cancel func()
	)
	defer ticker.Stop()
	for {
		select {
		case <-appctx.Done():
			return
		case <-ticker.C:
			appctx1 := appctx
			appctx1.Context, cancel = context.WithTimeout(appctx.Context, appctx.C.Timeout)
			Collect(appctx1)
			cancel()
		}
	}
}

func Collect(appctx Context) {
	var (
		tries     int
		timestamp = time.Now()
		err       error
	)
	for tries < appctx.C.Retries {
		tries++
		appctx.Log.Debugf("collecting records at timestamp %v. attempt %d", timestamp, tries)
		err = collect(appctx, timestamp)
		if err == nil {
			return
		}
		appctx.Log.Warnf("failed attempt %d, error %v", tries, err)
	}
	if err != nil {
		appctx.Log.Errorf("failed to collect records at timestamp %v. error %v", timestamp, err)
	}
	return
}

func collect(appctx Context, timestamp time.Time) error {
	req, err := http.NewRequestWithContext(appctx, "GET", appctx.C.URL, nil)
	if err != nil {
		return fmt.Errorf("creating http request failed: %w", err)
	}
	resp, err := appctx.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("http get failed: %w", err)
	}
	defer resp.Body.Close()

	appctx.Log.Debugf("got response using url %v. status %v.", appctx.C.URL, resp.Status)

	miners := &v1.MinerListResponse{}
	if err := jsonpb.Unmarshal(resp.Body, miners); err != nil {
		return fmt.Errorf("failed to unmarshal http response: %w", err)
	}

	appctx.Log.Debugf("total list of workers %d.", len(miners.Items))

	records := make([]interface{}, 0, len(miners.Items))
	for i := range miners.Items {
		if miners.Items[i].WorkerState != emitterv1.WorkerStateBonded {
			continue
		}
		records = append(records, Record{
			Timestamp: timestamp,
			Miner:     miners.Items[i],
		})
	}

	if err := appctx.DB.Save(appctx, records); err != nil {
		return fmt.Errorf("failed to save %d records to database: %w", len(records), err)
	}

	appctx.Log.Infof("saved bonded (%d) workers at timestamp %v", len(records), timestamp)
	return nil
}
