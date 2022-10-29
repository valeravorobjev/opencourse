package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// MgUser collection
type MgUser struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"` // User id
	Name       string             `bson:"name"`          // User name
	Avatar     string             `bson:"avatar"`        // Avatar image path
	Credential *MgCredential      `bson:"credential"`    // User credential properties
	Rating     int                `bson:"rating"`        // User rating
	Email      string             `bson:"email"`         // User email address
}

// MgCategory of curses collection
type MgCategory struct {
	Id            primitive.ObjectID `bson:"_id,omitempty"`            // Category id
	Lang          string             `bson:"lang"`                     // Language
	Name          string             `bson:"name"`                     // Category name
	LogoImg       string             `bson:"logo_img"`                 // Logo image icon
	SubCategories []*MgSubCategory   `bson:"sub_categories,omitempty"` // Sub categories
}

// MgStage collection
type MgStage struct {
	Id         primitive.ObjectID  `bson:"_id,omitempty"` // Stage id
	CourseId   primitive.ObjectID  `bson:"course_id"`     // Course id. One course has many stages
	Name       string              `bson:"name"`          // Course stage name
	HeaderImg  string              `bson:"header_img"`    // Stage header image
	Contents   []*MgPostContent    `bson:"contents"`      // Stage contents
	DateCreate primitive.Timestamp `bson:"date_create"`   // Create date of course
	DateUpdate primitive.Timestamp `bson:"date_update"`   // Update date of course
}

// MgCourse collection
type MgCourse struct {
	Id                primitive.ObjectID  `bson:"_id,omitempty"`       // Course id
	CategoryId        primitive.ObjectID  `bson:"category_id"`         // Course category
	SubCategoryNumber int                 `bson:"sub_category_number"` // Course sub category
	Lang              string              `bson:"lang"`                // Support languages
	Name              string              `bson:"name"`                // Course names
	Tags              []string            `bson:"tags,omitempty"`      // Course tags
	HeaderImg         string              `bson:"header_img"`          // Course header image
	DateCreate        primitive.Timestamp `bson:"date_create"`         // Create date of course
	DateUpdate        primitive.Timestamp `bson:"date_update"`         // Update date of course
	Rating            int                 `bson:"rating"`              // Course rating
	Description       string              `bson:"description"`         // Course description
	Actions           []*MgAction         `bson:"actions,omitempty"`   // Actions for comments
	Comments          []*MgComment        `bson:"comments,omitempty"`  // Comments for this post
}
