package mongodb

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"opencourse/common"
	"time"
)

//AddUser create user and save his to database
func (ctx *ContextMongoDb) AddUser(createUserQuery *common.OpenCreateUserQuery) (string, error) {
	col := ctx.client.Database(DbName).Collection(UserCollection)

	// Validate create user model

	if len(createUserQuery.Login) == 0 {
		return "", errors.New("user login is empty")
	}

	if len(createUserQuery.Password) == 0 || len(createUserQuery.Password) < 5 {
		return "", errors.New("password empty or has too small length, min length 5 symbols")
	}

	if len(createUserQuery.Name) == 0 {
		return "", errors.New("user name empty")
	}

	if len(createUserQuery.Roles) == 0 {
		return "", errors.New("roles is empty")
	}

	for _, role := range createUserQuery.Roles {
		switch role {
		case RoleUser, RoleAuthor, RoleAdmin:
			continue
		default:
			return "", errors.New("roles array contains non valid role")
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
		return "", nil
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), nil
}
