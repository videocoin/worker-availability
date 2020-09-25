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

type DB struct {
	Client *mongo.Collection
}

func (db *DB) Save(ctx context.Context, records []interface{}) error {
	_, err := db.Client.InsertMany(ctx, records, new(options.InsertManyOptions))
	return err
}

func (db *DB) All(ctx context.Context) ([]Record, error) {
	cursor, err := db.Client.Find(ctx, bson.M{}, new(options.FindOptions))
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
