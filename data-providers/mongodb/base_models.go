package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// Comment for post, course, etc..
type Comment struct {
	Id       primitive.ObjectID `bson:"id"`                  // Comment id
	UserId   primitive.ObjectID `bson:"user_id"`             // User id
	Text     string             `bson:"text"`                // Comment text
	Actions  []*Action          `bson:"actions,omitempty"`   // Actions for comment
	ParentId primitive.ObjectID `bson:"parent_id,omitempty"` // ParentId parent's comment id
}

// Action of user
type Action struct {
	UserId     primitive.ObjectID `bson:"user_id"`     // User id
	ActionType string             `bson:"action_type"` // Action type
}

// PostContent contains text and media data for stages
type PostContent struct {
	Body       string   `bson:"body"`                  // Post's body
	MediaItems []string `bson:"media_items,omitempty"` // Various attachments
}

// SubCategory for category
type SubCategory struct {
	Number int    `bson:"number"` // Sub category number
	Name   string `bson:"name"`   // Sub category name
}

// GlobStr contains localization text for some language. Check lang for get language type
type GlobStr struct {
	Lang string `bson:"lang"` // LangType is a type of text language.
	Text string `bson:"text"` // Text specific for language type
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
