package stats

import (
	"context"
	"time"

	v1 "github.com/videocoin/cloud-api/miners/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Record struct {
	Timestamp time.Time
	Miner     *v1.MinerResponse
}

type Aggregated struct {
	Timestamp time.Time `bson:"_id"`
	Records   []Record
}

type DB struct {
	Client *mongo.Collection
}

func (db *DB) Save(ctx context.Context, records []interface{}) error {
	_, err := db.Client.InsertMany(ctx, records, new(options.InsertManyOptions))
	return err
}

func (db *DB) All(ctx context.Context) ([]Record, error) {
	opts := new(options.FindOptions)
	opts.SetSort(bson.D{{"timestamp", 1}})
	cursor, err := db.Client.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	var rst []Record
	for cursor.Next(ctx) {
		record := Record{}
		if err := cursor.Decode(&record); err != nil {
			return nil, err
		}
		rst = append(rst, record)
	}
	return rst, nil
}

// Process iterates over records in the specified range, aggregating them by timestamp.
// Iterator function must return false if iteration must be stopped.
func (db *DB) Process(ctx context.Context, start, end time.Time, f func(Aggregated) bool) error {
	opts := new(options.AggregateOptions)
	opts = opts.SetAllowDiskUse(true)
	pipe := []bson.M{
		{
			"$match": bson.D{
				{"timestamp", bson.M{"$gte": start}},
				{"timestamp", bson.M{"$lte": end}}},
		},
		{
			"$group": bson.M{
				"_id":     "$timestamp",
				"records": bson.M{"$push": bson.M{"miner": "$miner"}},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}
	cursor, err := db.Client.Aggregate(ctx, pipe, opts)
	if err != nil {
		return err
	}
	for cursor.Next(ctx) {
		record := Aggregated{}
		if err := cursor.Decode(&record); err != nil {
			return err
		}
		if !f(record) {
			return nil
		}
	}
	return nil
}
