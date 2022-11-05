package database

// Collections names
const (
	UserCollection        = "users"         // Collection for store users
	CategoryCollection    = "categories"    // Collection for course categories
	StageCollection       = "stages"        // Collection store stages for courses
	CourseCollection      = "courses"       // Collection store courses
	TestCollection        = "tests"         // Collection for store stage's tests
	UserTestCollection    = "user_tests"    // Collection for store user and test relations and passed status
	UserConfirmCollection = "user_confirms" // Collection for store confirmation link for user registration. Use TTL index for auto remove documents.
)

const DbName = "opencourse" // Database name
