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
	Id        primitive.ObjectID `bson:"_id,omitempty"` // Category id
	Lang      string             `bson:"lang"`          // Language
	Name      string             `bson:"name"`          // Category name
	IconImg   string             `bson:"icon_img"`      // Icon for category
	HeaderImg string             `bson:"header_img"`    // Header image
}

// DbCourse collection
type DbCourse struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`  // Course id
	CategoryId  primitive.ObjectID `bson:"category_id"`    // Course category
	Enabled     bool               `bson:"enabled"`        // Enabled course
	Name        string             `bson:"name"`           // Course names
	Tags        []string           `bson:"tags,omitempty"` // Course tags
	Rating      int                `bson:"rating"`         // Course rating
	Description string             `bson:"description"`    // Course description
	IconImg     string             `bson:"icon_img"`       // Icon for category
	HeaderImg   string             `bson:"header_img"`     // Header image
}

// DbStage collection
type DbStage struct {
	Id        primitive.ObjectID `bson:"_id,omitempty"` // Stage id
	CourseId  primitive.ObjectID `bson:"course_id"`     // Course id. One course has many stages
	Name      string             `bson:"name"`          // Course stage name
	Contents  []*DbPostContent   `bson:"contents"`      // Stage contents
	HeaderImg string             `bson:"header_img"`    // Header image
}

// DbTest collection
type DbTest struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`          // Test id
	StageId       primitive.ObjectID `bson:"stage_id"`               // Stage id
	TestType      string             `bson:"test_type"`              // Test type
	LemmingsCount int                `bson:"lemmings_count"`         // Count of lemmings for passed test
	OptionTest    *DbOptionTest      `bson:"option_test,omitempty"`  // Option test. Test with option variant answers. Optional
	RewriteTest   *DbRewriteTest     `bson:"rewrite_test,omitempty"` // Rewrite test. Test with phrase how need write. Optional
}

// DbUserTest collection
type DbUserTest struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`     // User and Test records
	TestId   primitive.ObjectID `bson:"test_id,omitempty"` // Test id
	UserId   primitive.ObjectID `bson:"user_id"`           // User id
	IsPassed bool               `bson:"is_passed"`         // Check passed test
}
