package mongodb

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
AddUser create user and save his to database
@createUserQuery - create user model
*/
func (ctx *ContextMongoDb) AddUser(createUserQuery *common.OpenCreateUserQuery) (string, error) {
	col := ctx.Client.Database(DbName).Collection(UserCollection)

	// Validate create user model

	if len(createUserQuery.Login) == 0 {
		return "", openerrors.OpenFieldEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Login",
		}
	}

	if len(createUserQuery.Password) == 0 {
		return "", openerrors.OpenFieldEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Password",
		}
	}

	if len(createUserQuery.Password) < 5 {
		return "", openerrors.OpenMinLenErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field:  "createUserQuery.Password",
			MinLen: 5,
		}
	}

	if len(createUserQuery.Name) == 0 {
		return "", openerrors.OpenFieldEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Name",
		}
	}

	if len(createUserQuery.Roles) == 0 {
		return "", openerrors.OpenFieldEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/user_impl.go",
				Method: "AddUser",
			},
			Field: "createUserQuery.Roles",
		}
	}

	for _, role := range createUserQuery.Roles {
		switch role {
		case RoleUser, RoleAuthor, RoleAdmin:
			continue
		default:
			return "", openerrors.OpenRoleUnknownErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/user_impl.go",
					Method: "AddUser",
				},
				Role:  role,
				Roles: []string{RoleUser, RoleAuthor, RoleAdmin},
			}
		}
	}

	// Create new user

	var user User

	user.Name = createUserQuery.Name
	user.Avatar = createUserQuery.Avatar
	user.Email = createUserQuery.Email
	user.Rating = 0

	rand.Seed(time.Now().UnixNano())
	minRand := 10000000
	maxRand := 99999999
	salt := rand.Intn(maxRand-minRand) + minRand

	sha := sha1.New()
	str := fmt.Sprintf(createUserQuery.Password, salt)
	sha.Write([]byte(str))
	hash := hex.EncodeToString(sha.Sum(nil))

	timeNow := primitive.Timestamp{T: uint32(time.Now().Unix())}

	user.Credential = &Credential{
		Login:            createUserQuery.Login,
		Password:         hash,
		Salt:             salt,
		Roles:            createUserQuery.Roles,
		IsActive:         true,
		DateRegistration: timeNow,
		UpTime:           timeNow,
	}

	// Save user to DB

	result, err := col.InsertOne(context.Background(), user)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/user_impl.go",
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
