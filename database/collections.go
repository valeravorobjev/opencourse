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

// DbUserConfirm collection
type DbUserConfirm struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`    // User id
	ExpirationTime primitive.DateTime `bson:"expiration_time"`  // Expiration time for auto remove
	Login          string             `bson:"login"`            // User login
	Password       string             `bson:"password"`         // User password
	Name           string             `bson:"name"`             // User display name
	Email          string             `bson:"email"`            // Email user address
	Avatar         string             `bson:"avatar,omitempty"` // User avatar image path
	ConfirmaCode   string             `bson:"confirm_code"`     // Confirmation code for registration
	Confirmed      bool               `bson:"confirmed"`        // Confirmed if true
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
	DateCreate  primitive.DateTime `bson:"date_create"`    // Date create course
	DateUpdate  primitive.DateTime `bson:"date_update"`    // Date update course
}

// DbCoursePromotion collection
type DbCoursePromotion struct {
	Id             primitive.ObjectID `bson:"_id,omitempty"`   // Course promotion id
	CourseId       primitive.ObjectID `bson:"course_id"`       // Course id
	PromotionType  string             `bson:"promotion_type"`  // Promotion type
	Label          string             `bson:"label"`           // Promotion label text
	ExpirationTime primitive.DateTime `bson:"expiration_time"` // Promotion expiration time. After this time doc will be removed
}

// DbStage collection
type DbStage struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"` // Stage id
	CourseId    primitive.ObjectID `bson:"course_id"`     // Course id. One course has many stages
	Name        string             `bson:"name"`          // Course stage name
	Content     *DbPostContent     `bson:"content"`       // Stage contents
	HeaderImg   string             `bson:"header_img"`    // Header image
	OrderNumber int                `bson:"order_number"`  // Stage order number
}

// DbTest collection
type DbTest struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`          // Test id
	StageId       primitive.ObjectID `bson:"stage_id"`               // Stage id
	TestType      string             `bson:"test_type"`              // Test type
	LemmingsCount int                `bson:"lemmings_count"`         // Count of lemmings for passed test
	OptionTest    *DbOptionTest      `bson:"option_test,omitempty"`  // Option test. Test with option variant answers. Optional
	RewriteTest   *DbRewriteTest     `bson:"rewrite_test,omitempty"` // Rewrite test. Test with phrase how need write. Optional
	OrderNumber   int                `bson:"order_number"`           // Test order number
}

// DbUserTest collection
type DbUserTest struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`     // User and Test records
	TestId   primitive.ObjectID `bson:"test_id,omitempty"` // Test id
	UserId   primitive.ObjectID `bson:"user_id"`           // User id
	IsPassed bool               `bson:"is_passed"`         // Check passed test
}
