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
	OpenRoleAuthor = "Author"        // Is an author of courses
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

// OpenAddUserQuery model for create user
type OpenAddUserQuery struct {
	Login    string   `json:"login"`    // User login
	Password string   `json:"password"` // User password
	Email    string   `json:"email"`    // Email user address
	Name     string   `json:"name"`     // User display name
	Avatar   string   `json:"avatar"`   // User avatar image path
	Roles    []string `json:"roles"`    // User roles
}

// OpenGlobStr global text
type OpenGlobStr struct {
	Lang string `bson:"lang"` // LangType is a type of text language.
	Text string `bson:"text"` // Text specific for language type
}

// OpenSubCategory sub category
type OpenSubCategory struct {
	Number int            `json:"number"`
	Names  []*OpenGlobStr `json:"names"`
}

// OpenCategory category for course
type OpenCategory struct {
	Id            string             `json:"id"`
	Langs         []string           `json:"langs"` // Support languages
	Names         []*OpenGlobStr     `json:"names"`
	SubCategories []*OpenSubCategory `json:"sub_categories"`
}

// OpenAddCategoryQuery model for create category
type OpenAddCategoryQuery struct {
	Name          *OpenGlobStr       `json:"name"`
	SubCategories []*OpenSubCategory `json:"sub_categories"`
}

// OpenUpdateCategoryQuery model for update category
type OpenUpdateCategoryQuery struct {
	CategoryId    string             `json:"category_id"`
	Name          *OpenGlobStr       `json:"name"`
	SubCategories []*OpenSubCategory `json:"sub_categories"`
}

// OpenCourse model
type OpenCourse struct {
	Id                string         `json:"id"`                  // Course id
	AuthorIds         []string       `json:"author_ids"`          // Course authors ids
	Langs             []string       `json:"langs"`               // Support languages
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

// OpenComment user comment
type OpenComment struct {
	Id       string        `json:"id"`        // Comment id
	UserId   string        `json:"user_id"`   // User id
	Text     string        `json:"text"`      // Comment text
	Actions  []*OpenAction `json:"actions"`   // Actions for comment
	ParentId string        `json:"parent_id"` // ParentId parent's comment id
}

// OpenAction user action
type OpenAction struct {
	UserId     string `json:"user_id"`
	ActionType string `json:"action_type"`
}

// OpenAddCourseQuery add course query
type OpenAddCourseQuery struct {
	Langs             []string       `json:"langs"`               // Support languages
	Names             []*OpenGlobStr `json:"names"`               // Course name
	CategoryId        string         `json:"category_id"`         // Course category
	SubCategoryNumber int            `json:"sub_category_number"` // Course sub category
	Tags              []*OpenGlobStr `json:"tags"`                // Course tags
	HeaderImg         string         `json:"header_img"`          // Course header image
}
