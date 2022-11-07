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

// User roles
const (
	RoleUser  = "user"  // Simple user role
	RoleAdmin = "admin" // Privileged user role
)

// Test types
const (
	TestOption  = "option"  // Test with options (variant answers)
	TestRewrite = "rewrite" // Rewrite test. User need write key text from the question
)

// Promotion types
const (
	PromotionNew    = "new"    // New promotion record
	PromotionUpdate = "update" // Updated promotion record
)

// UserPreview user preview model for client
type UserPreview struct {
	Id     string   `json:"id"`     // User id
	Login  string   `json:"login"`  // User login
	Name   string   `json:"name"`   // User display name
	Email  string   `json:"email"`  // Email user address
	Avatar string   `json:"avatar"` // User avatar image path
	Rating int      `json:"rating"` // User rating
	Roles  []string `json:"roles"`  // User roles
}

// User model
type User struct {
	Id         string      `json:"id"`         // User id
	Name       string      `json:"name"`       // User name
	Avatar     string      `json:"avatar"`     // Avatar image path
	Credential *Credential `json:"credential"` // User credential properties
	Rating     int         `json:"rating"`     // User rating
	Email      string      `json:"email"`      // User email address
}

type Credential struct {
	Login            string    `json:"login"`             // User login
	Password         string    `json:"password"`          // User password
	Salt             int       `json:"salt"`              // Salt for generate password
	Roles            []string  `json:"roles"`             // User roles
	IsActive         bool      `json:"is_active"`         // Is user active or not
	DateRegistration time.Time `json:"date_registration"` // User registration date
	UpTime           time.Time `json:"uptime"`            // User uptime
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

// Category for course
type Category struct {
	Id        string `json:"id"`         // Category id
	Lang      string `json:"lang"`       // Support language
	Name      string `json:"name"`       // Category name
	IconImg   string `json:"icon_img"`   // Icon for category
	HeaderImg string `json:"header_img"` // Header image
}

// AddCategoryQuery model for create category
type AddCategoryQuery struct {
	Lang      string `json:"lang"`       // Support language
	Name      string `json:"name"`       // Category name
	IconImg   string `json:"icon_img"`   // Icon for category
	HeaderImg string `json:"header_img"` // Header image
}

// UpdateCategoryQuery model for update category
type UpdateCategoryQuery struct {
	CategoryId string `json:"category_id"` // Category id
	Lang       string `json:"lang"`        // Support language
	Name       string `json:"name"`        // Category name
}

// Course model
type Course struct {
	Id          string    `json:"id"`                    // Course id
	Name        string    `json:"name"`                  // Course name
	CategoryId  string    `json:"category_id"`           // Course category
	Enabled     bool      `json:"enabled"`               // Enabled course
	Tags        []string  `json:"tags"`                  // Course tags
	Rating      int       `json:"rating"`                // Course rating
	Description string    `json:"description,omitempty"` // Course description
	IconImg     string    `json:"icon_img"`              // Icon for category
	HeaderImg   string    `json:"header_img"`            // Header image
	DateCreate  time.Time `json:"date_create"`           // Date create course
	DateUpdate  time.Time `json:"date_update"`           // Date update course
}

// DbCoursePromotion collection
type DbCoursePromotion struct {
	Id             string    `json:"id"`              // Course promotion id
	CourseId       string    `json:"course_id"`       // Course id
	PromotionType  string    `json:"promotion_type"`  // Promotion type
	Label          string    `json:"label"`           // Promotion label text
	ExpirationTime time.Time `json:"expiration_time"` // Promotion expiration time. After this time doc will be removed
}

// AddCourseQuery add course query
type AddCourseQuery struct {
	Name        string   `json:"name"`        // Course name
	CategoryId  string   `json:"category_id"` // Course category
	Tags        []string `json:"tags"`        // Course tags
	Description string   `json:"description"` // Course description
	IconImg     string   `json:"icon_img"`    // Icon for category
	HeaderImg   string   `json:"header_img"`  // Header image
}

type PostContent struct {
	Body       string   `json:"body"`        // Post's body
	MediaItems []string `json:"media_items"` // Various attachments
}

type Stage struct {
	Id          string       `json:"id"`           // Stage id
	CourseId    string       `json:"course_id"`    // Course id. One course has many stages
	Name        string       `json:"name"`         // Course stage name
	Content     *PostContent `json:"content"`      // Stage contents
	HeaderImg   string       `json:"header_img"`   // Header image
	OrderNumber int          `json:"order_number"` // Stage order number
}

type StagePreview struct {
	Id          string `json:"id"`           // Stage id
	CourseId    string `json:"course_id"`    // Course id. One course has many stages
	Name        string `json:"name"`         // Course stage name
	HeaderImg   string `json:"header_img"`   // Header image
	OrderNumber int    `json:"order_number"` // Stage order number
}

type AddStageQuery struct {
	CourseId    string       `json:"course_id"`    // Course id. One course has many stages
	Name        string       `json:"name"`         // Course stage name
	Content     *PostContent `json:"content"`      // Stage contents
	HeaderImg   string       `json:"header_img"`   // Header image
	OrderNumber int          `json:"order_number"` // Stage order number
}

type UpdateStageQuery struct {
	StageId     string       `json:"stage_id"`     // Stage id
	CourseId    string       `json:"course_id"`    // Course id. One course has many stages
	Name        string       `json:"name"`         // Course stage name
	Content     *PostContent `json:"content"`      // Stage contents
	HeaderImg   string       `json:"header_img"`   // Header image
	OrderNumber int          `json:"order_number"` // Stage order number
}

type Test struct {
	Id            string       `json:"_id,omitempty"`          // Test id
	StageId       string       `json:"stage_id"`               // Stage id
	TestType      string       `json:"test_type"`              // Test type
	LemmingsCount int          `json:"lemmings_count"`         // Count of lemmings for passed test
	OptionTest    *OptionTest  `json:"option_test,omitempty"`  // Option test. Test with option variant answers. Optional
	RewriteTest   *RewriteTest `json:"rewrite_test,omitempty"` // Rewrite test. Test with phrase how need write. Optional
	OrderNumber   int          `json:"order_number"`           // Test order number
}

type TestPreview struct {
	Id            string `json:"_id,omitempty"`  // Test id
	StageId       string `json:"stage_id"`       // Stage id
	TestType      string `json:"test_type"`      // Test type
	LemmingsCount int    `json:"lemmings_count"` // Count of lemmings for passed test
	OrderNumber   int    `json:"order_number"`   // Test order number
}

type Option struct {
	Answer  string `json:"answer"`
	IsRight bool   `json:"is_right"`
}

type OptionTest struct {
	Question string    `json:"question"`
	Options  []*Option `json:"options"`
}

type RewriteTest struct {
	Question    string `json:"question"`
	RightAnswer string `json:"right_answer"`
}

type AddTestQuery struct {
	StageId       string       `json:"stage_id"`               // Stage id
	TestType      string       `json:"test_type"`              // Test type
	LemmingsCount int          `json:"lemmings_count"`         // Count of lemmings for passed test
	OptionTest    *OptionTest  `json:"option_test,omitempty"`  // Option test. Test with option variant answers. Optional
	RewriteTest   *RewriteTest `json:"rewrite_test,omitempty"` // Rewrite test. Test with phrase how need write. Optional
	OrderNumber   int          `json:"order_number"`           // Test order number
}

type LoginQuery struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterQuery struct {
	Login    string `json:"login"`    // User login
	Password string `json:"password"` // User password
	Name     string `json:"name"`     // User display name
	Email    string `json:"email"`    // Email user address
	Avatar   string `json:"avatar"`   // User avatar image path
}

type UserConfirm struct {
	Id             string    `json:"_id,omitempty"`    // User id
	ExpirationTime time.Time `json:"expiration_time"`  // Expiration time for auto remove
	Login          string    `json:"login"`            // User login
	Password       string    `json:"password"`         // User password
	Name           string    `json:"name"`             // User display name
	Email          string    `json:"email"`            // Email user address
	Avatar         string    `json:"avatar,omitempty"` // User avatar image path
	ConfirmaCode   string    `json:"confirm_code"`     // Confirmation code for registration
	Confirmed      bool      `json:"confirmed"`        // Confirmed if true
}
