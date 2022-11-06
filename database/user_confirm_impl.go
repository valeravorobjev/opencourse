package database

import (
	"context"
	"crypto/sha256"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"opencourse/common"
	"opencourse/common/openerrors"
	"time"
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

/*
GetUserConfirmByLogin return confirm data for user by login. Parameters:
login - user login;
*/
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

	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}

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

/*
AddUserConfirm add user confirm. Parameters:
query - RegisterQuery model;
*/
func (ctx *DbContext) AddUserConfirm(query *common.RegisterQuery) (*common.UserConfirm, error) {
	col := ctx.Client.Database(DbName).Collection(UserConfirmCollection)

	if query == nil {
		return nil, openerrors.ModelNilOrEmptyErr{
			Model: "addCategoryQuery",
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "AddUserConfirm",
			},
		}
	}

	var dbUserConfirm DbUserConfirm

	dbUserConfirm.ExpirationTime = primitive.NewDateTimeFromTime(time.Now().UTC().Add(time.Hour * 48))
	dbUserConfirm.Login = query.Login
	dbUserConfirm.Password = query.Password
	dbUserConfirm.Name = query.Name
	dbUserConfirm.Avatar = query.Avatar
	dbUserConfirm.Email = query.Email
	dbUserConfirm.Confirmed = false

	key := fmt.Sprintf("%s %s %s", dbUserConfirm.Login, dbUserConfirm.Email, dbUserConfirm.Password)
	h := sha256.New()
	h.Write([]byte(key))

	code := string(h.Sum(nil))

	dbUserConfirm.ConfirmaCode = code

	result, err := col.InsertOne(context.Background(), dbUserConfirm)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "AddUserConfirm",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	dbUserConfirm.Id = result.InsertedID.(primitive.ObjectID)

	userConfirm, err := dbUserConfirm.ToUserConfirm()

	if err != nil {
		return nil, openerrors.DefaultErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "AddUserConfirm",
			},
			Msg: err.Error(),
		}
	}

	return userConfirm, nil
}

/*
DeleteUserConfirm delete userConfirm. Parameters:
userConfirmId - user confirm id;
*/
func (ctx *DbContext) DeleteUserConfirm(userConfirmId string) error {
	col := ctx.Client.Database(DbName).Collection(UserConfirmCollection)

	objectUserConfirmId, err := primitive.ObjectIDFromHex(userConfirmId)

	if err != nil {
		return openerrors.InvalidIdErr{
			Id:        userConfirmId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/user_confirm_impl.go",
					Method: "DeleteUserConfirm",
				},
				Msg: err.Error(),
			},
		}
	}

	_, err = col.DeleteOne(context.Background(), bson.D{{"_id", objectUserConfirmId}})

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/user_confirm_impl.go",
				Method: "DeleteUserConfirm",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	return nil
}
