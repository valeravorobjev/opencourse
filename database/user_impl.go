package database

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"opencourse/common"
	"opencourse/common/openerrors"
	"time"
)

func BuildHash(password string, salt int) string {
	sha := sha1.New()
	str := fmt.Sprintf(password, salt)
	sha.Write([]byte(str))
	hash := hex.EncodeToString(sha.Sum(nil))

	return hash
}

/*
AddUser create user and save his to database. Parameters:
createUserQuery - create user model;
*/
func (ctx *DbContext) AddUser(createUserQuery *common.AddUserQuery) (string, error) {
	col := ctx.Client.Database(DbName).Collection(UserCollection)

	// Validate create user model

	if len(createUserQuery.Login) == 0 {
		return "", openerrors.FieldEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Login",
		}
	}

	if len(createUserQuery.Password) == 0 {
		return "", openerrors.FieldEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Password",
		}
	}

	if len(createUserQuery.Password) < 5 {
		return "", openerrors.MinLenErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field:  "createUserQuery.Password",
			MinLen: 5,
		}
	}

	if len(createUserQuery.Name) == 0 {
		return "", openerrors.FieldEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Name",
		}
	}

	if len(createUserQuery.Roles) == 0 {
		return "", openerrors.FieldEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Roles",
		}
	}

	for _, role := range createUserQuery.Roles {
		switch role {
		case common.RoleUser, common.RoleAuthor, common.RoleAdmin:
			continue
		default:
			return "", openerrors.RoleUnknownErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/mongodb/user_impl.go",
					Method: "AddUser",
				},
				Role:  role,
				Roles: []string{common.RoleUser, common.RoleAuthor, common.RoleAdmin},
			}
		}
	}

	// Create new user

	var dbUser DbUser

	dbUser.Name = createUserQuery.Name
	dbUser.Avatar = createUserQuery.Avatar
	dbUser.Email = createUserQuery.Email
	dbUser.Rating = 0

	rand.Seed(time.Now().UnixNano())
	minRand := 10000000
	maxRand := 99999999
	salt := rand.Intn(maxRand-minRand) + minRand

	hash := BuildHash(createUserQuery.Password, salt)

	timeNow := primitive.NewDateTimeFromTime(time.Now().UTC())

	dbUser.Credential = &DbCredential{
		Login:            createUserQuery.Login,
		Password:         hash,
		Salt:             salt,
		Roles:            createUserQuery.Roles,
		IsActive:         true,
		DateRegistration: timeNow,
		UpTime:           timeNow,
	}

	// Save user to DB

	result, err := col.InsertOne(context.Background(), dbUser)

	if err != nil {
		return "", openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "AddUser",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}

/*
GetUserByLogin return user by login. Parameters:
login - user login;
*/
func (ctx *DbContext) GetUserByLogin(login string) (*common.User, error) {
	col := ctx.Client.Database(DbName).Collection(UserCollection)

	if len(login) < 1 {
		return nil, openerrors.FieldEmptyErr{
			Field: "login",
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "GetUserByLogin",
			},
		}
	}

	find := bson.D{{
		"credential.login", login,
	}}

	var dbUser DbUser
	err := col.FindOne(context.Background(), find).Decode(&dbUser)

	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "GetUserByLogin",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	user, err := dbUser.ToUser()

	if err != nil {
		return nil, openerrors.DefaultErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/user_impl.go",
				Method: "GetUserByLogin",
			},
			Msg: err.Error(),
		}
	}

	return user, nil
}
