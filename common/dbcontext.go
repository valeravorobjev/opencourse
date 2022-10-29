package common

type DbContext interface {
	Defaults()
	Connect(uri string) error
	Disconnect() error

	AddUser(createUserQuery *AddUserQuery) (string, error)

	GetCategories(langs []string) ([]*OpenCategory, error)
	AddCategory(addCategoryQuery *AddCategoryQuery) (string, error)
	UpdateCategory(cid string, name string) error
	UpdateSubCategory(cid string, scn int, name string) error
	AddSubCategory(cid string, name string) error
	DeleteSubCategory(cid string, scn int) error

	ClearCourses() error
	GetCourse(id string) (*Course, error)
	GetCourses(take int64, skip int64) ([]*Course, error)
	AddCourse(userId string, addCourseQuery *AddCourseQuery) (string, error)
	AddCourseAction(id string, userId string, actionType string) error
	ChangeCourseAction(id string, userId string, actionType string) error
	RemoveCourseAction(id string, userId string) error
	AddCourseComment(id string, userId string, text string) (string, error)
	ReplyCourseComment(id string, userId string, commentId string, text string) error
	RemoveCourseComment(id string, commentId string) error
	AddCourseTags(id string, tags []string) error
	RemoveCourseTags(id string, tags []string) error
}
