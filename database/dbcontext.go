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
	ConStr          string        // Connection string
	DbName          string        // Db name. Example: mongodb/opencourse
	SmtpAccount     string        // SMTP account
	SmtpAccountPass string        // SMTP account password
	Endpoint        string        // Endpoint (base url)
	Client          *mongo.Client // Client connection for db
}

// Defaults init values
func (ctx *DbContext) Defaults(conStr string, smtpAccount string, smtpPass string, endpoint string) {

	ctx.DbName = fmt.Sprintf("mongodb/%s", DbName)
	ctx.SmtpAccountPass = smtpPass
	ctx.SmtpAccount = smtpAccount
	ctx.Endpoint = endpoint
	ctx.ConStr = conStr
}

// Connect to db
func (ctx *DbContext) Connect() error {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(ctx.ConStr))
	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/dbcontext.go",
				Method: "Connect",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
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
			ConStr: ctx.ConStr,
		}
	}

	return nil
}
