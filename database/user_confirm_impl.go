package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"opencourse/common"
	"opencourse/common/openerrors"
)

// ClearUserConfirms remove all data from user_confirms collection
func (ctx *DbContext) ClearUserConfirms() error {
	col := ctx.Client.Database(DbName).Collection(UserConfirmCollection)

	_, err := col.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "ClearUserConfirms",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
GetUserConfirm return user confirmation data. Parameters:
userConfirmId - user confirm id;
*/
func (ctx *DbContext) GetUserConfirm(userConfirmId string) (*common.UserConfirm, error) {

	col := ctx.Client.Database(DbName).Collection(UserConfirmCollection)

	objectUserConfirmId, err := primitive.ObjectIDFromHex(userConfirmId)

	if err != nil {
		return nil, openerrors.InvalidIdErr{
			Id:        userConfirmId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/user_confirm_impl.go",
					Method: "GetUserConfirm",
				},
				Msg: err.Error(),
			},
		}
	}

	filter := bson.D{{"_id", objectUserConfirmId}}

	var dbUserConfirm DbUserConfirm
	err = col.FindOne(context.Background(), filter).Decode(&dbUserConfirm)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "ClearUserConfirms",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	userConfirm, err := dbUserConfirm.ToUserConfirm()

	if err != nil {
		return nil, openerrors.DefaultErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "ClearUserConfirms",
			},
			Msg: err.Error(),
		}
	}

	return userConfirm, nil
}

/*
SetConfirmed set confirm for user
userConfirmId - user confirm id;
*/
func (ctx *DbContext) SetConfirmed(userConfirmId string) error {
	col := ctx.Client.Database(DbName).Collection(UserConfirmCollection)

	objectUserConfirmId, err := primitive.ObjectIDFromHex(userConfirmId)

	if err != nil {
		return openerrors.InvalidIdErr{
			Id:        userConfirmId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/user_confirm_impl.go",
					Method: "SetConfirmed",
				},
				Msg: err.Error(),
			},
		}
	}

	find := bson.D{{"_id", objectUserConfirmId}}
	update := bson.D{{
		"$set", bson.D{{"confirmed", true}},
	}}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "SetConfirmed",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	return nil
}

func (ctx *DbContext) GetUserConfirmByLogin(login string) (*common.UserConfirm, error) {
	col := ctx.Client.Database(DbName).Collection(UserConfirmCollection)

	if len(login) < 1 {
		return nil, openerrors.FieldEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "SetConfirmed",
			},
			Field: "login",
		}
	}

	find := bson.D{{"login", login}}

	var dbUserConfirm DbUserConfirm
	err := col.FindOne(context.Background(), find).Decode(&dbUserConfirm)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "GetUserConfirmByLogin",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	userConfirm, err := dbUserConfirm.ToUserConfirm()

	if err != nil {
		return nil, openerrors.DefaultErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "GetUserConfirmByLogin",
			},
			Msg: err.Error(),
		}
	}

	return userConfirm, nil

}
