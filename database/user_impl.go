package database

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"opencourse/common"
	"opencourse/common/openerrors"
	"time"
)

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

	sha := sha1.New()
	str := fmt.Sprintf(createUserQuery.Password, salt)
	sha.Write([]byte(str))
	hash := hex.EncodeToString(sha.Sum(nil))

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
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}
