package mongoStore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	connectTimeOut = time.Second * 30
)

func NewClient(ctx context.Context, mongoUrl, name string) (*mongo.Database, error) {
	clientOpt := options.Client().ApplyURI(mongoUrl)
	clientOpt.Timeout = &connectTimeOut

	client, err := mongo.Connect(ctx, clientOpt)

	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return client.Database(name), nil
}
