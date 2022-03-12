package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Database *mongo.Client
}

func NewMongoDB(addr string) (m *MongoDB, err error) {
	m = &MongoDB{}
	if m.Database, err = mongo.NewClient(options.Client().ApplyURI(addr)); err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = m.Database.Connect(ctx)
	return
}
