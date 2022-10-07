package mongodb

// Collections names
const (
	UserCollection     = "users"      // Collection for store users
	CategoryCollection = "categories" // Collection for course categories
	StageCollection    = "stages"     // Collection store stages for courses
	CourseCollection   = "courses"    // Collection store courses
)

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
	RoleUser   string = "User"          // Simple user role
	RoleAuthor        = "Author"        // Is a author of courses
	RoleAdmin         = "Administrator" // Privileged user role
)

const DbName = "opencourse" // Database name
