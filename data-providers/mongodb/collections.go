package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// User collections
type User struct {
	Id         primitive.ObjectID `bson:"_id"`        // User id
	Name       string             `bson:"name"`       // User name
	Credential Credential         `bson:"credential"` // User credential properties
	Rating     int                `bson:"rating"`     // User rating
}

// Course collection
type Course struct {
	Id         primitive.ObjectID  `bson:"_id"`         // Course id
	AuthorId   primitive.ObjectID  `bson:"author_id"`   // Course author id
	Title      string              `bson:"title"`       // Course title
	TitleImage string              `bson:"title_image"` // Course header image
	PostBody   PostBody            `bson:"post_body"`   // Body of course
	DateCreate primitive.Timestamp `bson:"date_create"` // Create date of course
	DateUpdate primitive.Timestamp `bson:"date_update"` // Update date of course
	Rating     int                 `bson:"rating"`      // Course rating
	Actions    []*Action           `bson:"actions"`     // Actions for comments
	Comments   []*Comment          `bson:"comments"`    // Comments for this post
}
