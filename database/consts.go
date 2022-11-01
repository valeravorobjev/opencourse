package database

// Collections names
const (
	UserCollection     = "users"      // Collection for store users
	CategoryCollection = "categories" // Collection for course categories
	StageCollection    = "stages"     // Collection store stages for courses
	CourseCollection   = "courses"    // Collection store courses
	TestsCollection    = "tests"      // Collection for store stage's tests
	UserTests          = "user_tests" // Collection for store user and test relations and passed status
)

const DbName = "opencourse" // Database name
