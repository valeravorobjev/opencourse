package common

import (
	"time"
)

// Open - project prefix

// Languages types
const (
	LangEn = "en" // English
	LangFr = "fr" // France
	LangDe = "de" // Germany
	LangIt = "it" // Italian
	LangRu = "ru" // Russian
)

// Action types
const (
	ActionLike    = "Like"    // Like
	ActionDislike = "Dislike" // Dislike
)

// User roles
const (
	RoleUser   = "User"          // Simple user role
	RoleAuthor = "Author"        // Is an author of courses
	RoleAdmin  = "Administrator" // Privileged user role
)

// User user common model
type User struct {
	Id     string   `json:"id"`     // User id
	Login  string   `json:"login"`  // User login
	Name   string   `json:"name"`   // User display name
	Email  string   `json:"email"`  // Email user address
	Avatar string   `json:"avatar"` // User avatar image path
	Rating int      `json:"rating"` // User rating
	Roles  []string `json:"roles"`  // User roles
}

// AddUserQuery model for create user
type AddUserQuery struct {
	Login    string   `json:"login"`    // User login
	Password string   `json:"password"` // User password
	Email    string   `json:"email"`    // Email user address
	Name     string   `json:"name"`     // User display name
	Avatar   string   `json:"avatar"`   // User avatar image path
	Roles    []string `json:"roles"`    // User roles
}

// SubCategory sub category
type SubCategory struct {
	Number int    `json:"number"`
	Name   string `json:"name"`
}

// Category for course
type Category struct {
	Id            string         `json:"id"`
	Lang          string         `json:"lang"` // Support language
	Name          string         `json:"name"`
	SubCategories []*SubCategory `json:"sub_categories"`
}

// AddCategoryQuery model for create category
type AddCategoryQuery struct {
	Lang          string         `json:"lang"` // Support language
	Name          string         `json:"name"`
	SubCategories []*SubCategory `json:"sub_categories"`
}

// UpdateCategoryQuery model for update category
type UpdateCategoryQuery struct {
	CategoryId    string         `json:"category_id"`
	Langs         []string       `json:"lang"` // Support language
	Name          string         `json:"name"`
	SubCategories []*SubCategory `json:"sub_categories"`
}

// Course model
type Course struct {
	Id                string     `json:"id"`                    // Course id
	Lang              string     `json:"lang"`                  // Support language
	Name              string     `json:"name"`                  // Course name
	CategoryId        string     `json:"category_id"`           // Course category
	SubCategoryNumber int        `json:"sub_category_number"`   // Course sub category
	Tags              []string   `json:"tags"`                  // Course tags
	HeaderImg         string     `json:"header_img"`            // Course header image
	DateCreate        time.Time  `json:"date_create"`           // Create date of course
	DateUpdate        time.Time  `json:"date_update"`           // Update date of course
	Rating            int        `json:"rating"`                // Course rating
	Description       string     `json:"description,omitempty"` // Course description
	Actions           []*Action  `json:"actions,omitempty"`     // Actions for comments
	Comments          []*Comment `json:"comments,omitempty"`    // Comments for this post
}

// Comment user comment
type Comment struct {
	Id       string    `json:"id"`        // Comment id
	UserId   string    `json:"user_id"`   // User id
	Text     string    `json:"text"`      // Comment text
	Actions  []*Action `json:"actions"`   // Actions for comment
	ParentId string    `json:"parent_id"` // ParentId parent's comment id
}

// Action user action
type Action struct {
	UserId     string `json:"user_id"`
	ActionType string `json:"action_type"`
}

// AddCourseQuery add course query
type AddCourseQuery struct {
	Lang              string   `json:"lang"`                // Support languages
	Name              string   `json:"name"`                // Course name
	CategoryId        string   `json:"category_id"`         // Course category
	SubCategoryNumber int      `json:"sub_category_number"` // Course sub category
	Tags              []string `json:"tags"`                // Course tags
	HeaderImg         string   `json:"header_img"`          // Course header image
}
