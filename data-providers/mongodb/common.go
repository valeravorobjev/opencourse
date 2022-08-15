package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// Comment for post, course, etc..
type Comment struct {
	UserId   primitive.ObjectID `bson:"user_id"`  // User id
	Text     string             `bson:"text"`     // Comment text
	Actions  []*Action          `bson:"actions"`  // Actions for comment
	Comments []*Comment         `bson:"comments"` // Comments for this comment
}

// Action of user
type Action struct {
	UserId     primitive.ObjectID `bson:"user_id"`     // User id
	ActionType string             `bson:"action_type"` // Action type
}

// PostBody contains text and media data for stages
type PostBody struct {
	Body       *GlobalStr `bson:"body"`        // Post body
	MediaItems []string   `bson:"media_items"` // Various attachments
}

// SubCategory for category
type SubCategory struct {
	Number int        `bson:"number"` // Sub category number
	Name   *GlobalStr `bson:"name"`   // Sub category name
}

// Action types
const (
	ActionLike    = "Like"    // Like
	ActionDislike = "Dislike" // Dislike
)

// User roles
const (
	RoleUser  string = "User"          // Simple user role
	Author           = "Author"        // Is a author of courses
	RoleAdmin        = "Administrator" // Privileged user role
)

// GlobalStr contains localization text for some languages. Check types array for get language type
type GlobalStr struct {
	Types []string `bson:"types"` // Types of text context. Contains elements with no empty str.
	En    string   `bson:"en"`    // English text content
	It    string   `bson:"it"`    // Italian text content
	Fr    string   `bson:"fr"`    // France text content
	Ru    string   `bson:"ru"`    // Russian text content
}

// Credential contains user auth data
type Credential struct {
	Login            string              `bson:"login"`             // User login
	Password         string              `bson:"password"`          // User password
	Salt             int                 `bson:"salt"`              // Salt for generate password
	Roles            []string            `bson:"roles"`             // User roles
	IsActive         bool                `bson:"is_active"`         // Is user active or not
	DateRegistration primitive.Timestamp `bson:"date_registration"` // User registration date
	UpTime           primitive.Timestamp `bson:"uptime"`            // User uptime
}
