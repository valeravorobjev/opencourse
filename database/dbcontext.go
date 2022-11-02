package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common/openerrors"
)

// DbContext is a context for work with mongo db
type DbContext struct {
	Uri    string        // Connection string
	DbName string        // Db name. Example: mongodb/opencourse
	Client *mongo.Client // Client connection for db
}

// Defaults init values
func (ctx *DbContext) Defaults() {
	ctx.DbName = fmt.Sprintf("mongodb/%s", DbName)
}

// Connect to db
func (ctx *DbContext) Connect(uri string) error {
	ctx.Uri = uri
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/dbcontext.go",
				Method: "Connect",
			},
			DbName: ctx.DbName,
			ConStr: uri,
		}
	}
	ctx.Client = client
	return nil
}

// Disconnect db
func (ctx *DbContext) Disconnect() error {
	err := ctx.Client.Disconnect(context.Background())
	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/dbcontext.go",
				Method: "Disconnect",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
		}
	}

	return nil
}
