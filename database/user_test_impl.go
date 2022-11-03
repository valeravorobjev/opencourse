package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"opencourse/common/openerrors"
)

// ClearUserTests remove all data from user_tests collection
func (ctx *DbContext) ClearUserTests() error {
	col := ctx.Client.Database(DbName).Collection(UserTestCollection)

	_, err := col.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_test_impl.go",
				Method: "ClearUserTests",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
