package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// User collection
type User struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"` // User id
	Name       string             `bson:"name"`          // User name
	Avatar     string             `bson:"avatar"`        // Avatar image path
	Credential *Credential        `bson:"credential"`    // User credential properties
	Rating     int                `bson:"rating"`        // User rating
	Email      string             `bson:"email"`         // User email address
}

// Category of curses collection
type Category struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`            // Category id
	Langs         []string           `bson:"langs"`                    // Support languages
	Names         []*GlobStr         `bson:"names"`                    // Category names
	SubCategories []*SubCategory     `bson:"sub_categories,omitempty"` // Sub categories
}

// Stage collection
type Stage struct {
	Id         primitive.ObjectID  `bson:"_id,omitempty"`      // Stage id
	AuthorId   primitive.ObjectID  `bson:"author_id"`          // Author id how right stage
	CourseId   primitive.ObjectID  `bson:"course_id"`          // Course id. One course has many stages
	Names      []*GlobStr          `bson:"names"`              // Course stage names
	HeaderImg  string              `bson:"header_img"`         // Stage header image
	Contents   []*PostContent      `bson:"contents"`           // Stage contents
	DateCreate primitive.Timestamp `bson:"date_create"`        // Create date of course
	DateUpdate primitive.Timestamp `bson:"date_update"`        // Update date of course
	Actions    []*Action           `bson:"actions,omitempty"`  // Actions for comments
	Comments   []*Comment          `bson:"comments,omitempty"` // Comments for this post
}

// Course collection
type Course struct {
	Id                primitive.ObjectID   `bson:"_id,omitempty"`       // Course id
	AuthorIds         []primitive.ObjectID `bson:"author_ids"`          // Course author ids
	Langs             []string             `bson:"langs"`               // Support languages
	Names             []*GlobStr           `bson:"names"`               // Course names
	CategoryId        primitive.ObjectID   `bson:"category_id"`         // Course category
	SubCategoryNumber int                  `bson:"sub_category_number"` // Course sub category
	Tags              []*GlobStr           `bson:"tags,omitempty"`      // Course tags
	HeaderImg         string               `bson:"header_img"`          // Course header image
	DateCreate        primitive.Timestamp  `bson:"date_create"`         // Create date of course
	DateUpdate        primitive.Timestamp  `bson:"date_update"`         // Update date of course
	Rating            int                  `bson:"rating"`              // Course rating
	Actions           []*Action            `bson:"actions,omitempty"`   // Actions for comments
	Comments          []*Comment           `bson:"comments,omitempty"`  // Comments for this post
}
