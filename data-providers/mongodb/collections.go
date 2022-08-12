package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// User collection
type User struct {
	Id         primitive.ObjectID `bson:"_id"`        // User id
	Name       string             `bson:"name"`       // User name
	Credential *Credential        `bson:"credential"` // User credential properties
	Rating     int                `bson:"rating"`     // User rating
}

// Stage collection
type Stage struct {
	Id         primitive.ObjectID  `bson:"_id"`         // Stage id
	CourseId   primitive.ObjectID  `bson:"course_id"`   // Course id
	Name       string              `bson:"name"`        // Course stage name
	BannerImg  string              `bson:"banner_img"`  // Stage header image
	PostBody   *PostBody           `bson:"post_body"`   // Body of stage
	DateCreate primitive.Timestamp `bson:"date_create"` // Create date of course
	DateUpdate primitive.Timestamp `bson:"date_update"` // Update date of course
	Actions    []*Action           `bson:"actions"`     // Actions for comments
	Comments   []*Comment          `bson:"comments"`    // Comments for this post
}

// Course collection
type Course struct {
	Id          primitive.ObjectID  `bson:"_id"`          // Course id
	AuthorId    primitive.ObjectID  `bson:"author_id"`    // Course author id
	Name        string              `bson:"title"`        // Course title
	Category    string              `bson:"category"`     // Course category
	SubCategory string              `bson:"sub_category"` // Course sub category
	Tags        []string            `bson:"tags"`         // Course tags
	BannerImg   string              `bson:"banner_img"`   // Course header image
	DateCreate  primitive.Timestamp `bson:"date_create"`  // Create date of course
	DateUpdate  primitive.Timestamp `bson:"date_update"`  // Update date of course
	Rating      int                 `bson:"rating"`       // Course rating
	Actions     []*Action           `bson:"actions"`      // Actions for comments
	Comments    []*Comment          `bson:"comments"`     // Comments for this post
}
