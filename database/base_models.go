package database

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
This file contains sub documents for collections.
*/

// DbPostContent contains text and media data for stages
type DbPostContent struct {
	Body       string   `bson:"body"`                  // Post's body
	MediaItems []string `bson:"media_items,omitempty"` // Various attachments
}

// DbCredential contains user auth data
type DbCredential struct {
	Login            string             `bson:"login"`             // User login
	Password         string             `bson:"password"`          // User password
	Salt             int                `bson:"salt"`              // Salt for generate password
	Roles            []string           `bson:"roles"`             // User roles
	IsActive         bool               `bson:"is_active"`         // Is user active or not
	DateRegistration primitive.DateTime `bson:"date_registration"` // User registration date
	UpTime           primitive.DateTime `bson:"uptime"`            // User uptime
}

type Option struct {
	Answer  string `bson:"answer"`
	IsRight bool   `bson:"is_right"`
}

type DbOptionTest struct {
	Question string    `bson:"question"`
	Options  []*Option `bson:"options"`
}

type DbRewriteTest struct {
	Question    string `bson:"question"`
	RightAnswer bool   `bson:"right_answer"`
}
