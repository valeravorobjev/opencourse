package common

import (
	"time"
)

// Open - project prefix

// Languages types
const (
	OpenLangEn = "en" // English
	OpenLangFr = "fr" // France
	OpenLangDe = "de" // Germany
	OpenLangIt = "it" // Italian
	OpenLangRu = "ru" // Russian
)

// Action types
const (
	OpenActionLike    = "Like"    // Like
	OpenActionDislike = "Dislike" // Dislike
)

// User roles
const (
	OpenRoleUser   = "User"          // Simple user role
	OpenRoleAuthor = "Author"        // Is a author of courses
	OpenRoleAdmin  = "Administrator" // Privileged user role
)

// OpenUser user common model
type OpenUser struct {
	Id     string   `json:"id"`     // User id
	Login  string   `json:"login"`  // User login
	Name   string   `json:"name"`   // User display name
	Email  string   `json:"email"`  // Email user address
	Avatar string   `json:"avatar"` // User avatar image path
	Rating int      `json:"rating"` // User rating
	Roles  []string `json:"roles"`  // User roles
}

// OpenCreateUserQuery model for create user
type OpenAddUserQuery struct {
	Login    string   `json:"login"`    // User login
	Password string   `json:"password"` // User password
	Email    string   `json:"email"`    // Email user address
	Name     string   `json:"name"`     // User display name
	Avatar   string   `json:"avatar"`   // User avatar image path
	Roles    []string `json:"roles"`    // User roles
}

type OpenGlobStr struct {
	Lang string `bson:"lang"` // LangType is a type of text language.
	Text string `bson:"text"` // Text specific for language type
}

type OpenSubCategory struct {
	Number int          `json:"number"`
	Name   *OpenGlobStr `json:"name"`
}

type OpenCategory struct {
	Id            string             `json:"id"`
	Name          *OpenGlobStr       `json:"name"`
	SubCategories []*OpenSubCategory `json:"sub_categories"`
}

type OpenAddCategoryQuery struct {
	Name          *OpenGlobStr       `json:"name"`
	SubCategories []*OpenSubCategory `json:"sub_categories"`
}

type OpenUpdateCategoryQuery struct {
	CategoryId    string             `json:"category_id"`
	Name          *OpenGlobStr       `json:"name"`
	SubCategories []*OpenSubCategory `json:"sub_categories"`
}

// OpenCourse model
type OpenCourse struct {
	Id                string         `json:"id"`                  // Course id
	AuthorIds         []string       `json:"author_ids"`          // Course authors ids
	Names             []*OpenGlobStr `json:"names"`               // Course name
	CategoryId        string         `json:"category_id"`         // Course category
	SubCategoryNumber int            `json:"sub_category_number"` // Course sub category
	Tags              []*OpenGlobStr `json:"tags"`                // Course tags
	HeaderImg         string         `json:"header_img"`          // Course header image
	DateCreate        time.Time      `json:"date_create"`         // Create date of course
	DateUpdate        time.Time      `json:"date_update"`         // Update date of course
	Rating            int            `json:"rating"`              // Course rating
	Actions           []*OpenAction  `json:"actions"`             // Actions for comments
	Comments          []*OpenComment `json:"comments"`            // Comments for this post
}

type OpenComment struct {
	Id       string        `json:"id"`        // Comment id
	UserId   string        `json:"user_id"`   // User id
	Text     string        `json:"text"`      // Comment text
	Actions  []*OpenAction `json:"actions"`   // Actions for comment
	ParentId string        `json:"parent_id"` // ParentId parent's comment id
}

type OpenAction struct {
	UserId     string `json:"user_id"`
	ActionType string `json:"action_type"`
}

type OpenAddCourseQuery struct {
	Name              *OpenGlobStr   `json:"names"`               // Course name
	CategoryId        string         `json:"category_id"`         // Course category
	SubCategoryNumber int            `json:"sub_category_number"` // Course sub category
	Tags              []*OpenGlobStr `json:"tags"`                // Course tags
	HeaderImg         string         `json:"header_img"`          // Course header image
}
