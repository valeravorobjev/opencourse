package mongodb

import "go.mongodb.org/mongo-driver/bson/primitive"

// MgComment for post, course, etc..
type MgComment struct {
	Id       primitive.ObjectID `bson:"id"`                  // Comment id
	UserId   primitive.ObjectID `bson:"user_id"`             // User id
	Text     string             `bson:"text"`                // Comment text
	Actions  []*MgAction        `bson:"actions,omitempty"`   // Actions for comment
	ParentId primitive.ObjectID `bson:"parent_id,omitempty"` // ParentId parent's comment id
}

// MgAction of user
type MgAction struct {
	UserId     primitive.ObjectID `bson:"user_id"`     // User id
	ActionType string             `bson:"action_type"` // Action type
}

// MgPostContent contains text and media data for stages
type MgPostContent struct {
	Body       string   `bson:"body"`                  // Post's body
	MediaItems []string `bson:"media_items,omitempty"` // Various attachments
}

// MgSubCategory for category
type MgSubCategory struct {
	Number  int    `bson:"number"`   // Sub category number
	LogoImg string `bson:"logo_img"` // Logo image icon
	Name    string `bson:"name"`     // Sub category name
}

// MgCredential contains user auth data
type MgCredential struct {
	Login            string              `bson:"login"`             // User login
	Password         string              `bson:"password"`          // User password
	Salt             int                 `bson:"salt"`              // Salt for generate password
	Roles            []string            `bson:"roles"`             // User roles
	IsActive         bool                `bson:"is_active"`         // Is user active or not
	DateRegistration primitive.Timestamp `bson:"date_registration"` // User registration date
	UpTime           primitive.Timestamp `bson:"uptime"`            // User uptime
}
