package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// GetCategories return all categories from db
func (ctx *ContextMongoDb) GetCategories() ([]*Category, error) {
	col := ctx.client.Database(DbName).Collection(CategoryCollection)

	cursor, err := col.Find(context.Background(), bson.D{})

	if err != nil {
		return nil, err
	}

	var categories []*Category

	err = cursor.Decode(categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (ctx *ContextMongoDb) AddCategory(category *Category) (string, error) {
	col := ctx.client.Database(DbName).Collection(CategoryCollection)

	result, err := col.InsertOne(context.TODO(), category)

	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), err
}
