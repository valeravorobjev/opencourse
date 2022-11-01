package database

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
This file contains collection's structs.
MongoDB auto create same collections with these structs.
*/

// DbUser collection
type DbUser struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"` // User id
	Name       string             `bson:"name"`          // User name
	Avatar     string             `bson:"avatar"`        // Avatar image path
	Credential *DbCredential      `bson:"credential"`    // User credential properties
	Rating     int                `bson:"rating"`        // User rating
	Email      string             `bson:"email"`         // User email address
}

// DbCategory of curses collection
type DbCategory struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"` // Category id
	Lang string             `bson:"lang"`          // Language
	Name string             `bson:"name"`          // Category name
}

// DbCourse collection
type DbCourse struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`  // Course id
	CategoryId  primitive.ObjectID `bson:"category_id"`    // Course category
	Name        string             `bson:"name"`           // Course names
	Tags        []string           `bson:"tags,omitempty"` // Course tags
	Rating      int                `bson:"rating"`         // Course rating
	Description string             `bson:"description"`    // Course description
}

// DbStage collection
type DbStage struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"` // Stage id
	CourseId primitive.ObjectID `bson:"course_id"`     // Course id. One course has many stages
	Name     string             `bson:"name"`          // Course stage name
	Contents []*DbPostContent   `bson:"contents"`      // Stage contents
}

type DbTest[T any] struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	TestType string             `bson:"test_type"`
}
