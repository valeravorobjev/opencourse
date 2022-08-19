package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ContextMongoDb is a context for work with mongo db
type ContextMongoDb struct {
	client *mongo.Client // Client connection for db
}

// Connect to db
func (ctx *ContextMongoDb) Connect(uri string) error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}
	ctx.client = client
	return nil
}

// Disconnect db
func (ctx *ContextMongoDb) Disconnect() error {
	err := ctx.client.Disconnect(context.Background())
	if err != nil {
		return err
	}

	return nil
}
