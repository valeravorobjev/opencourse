package integration

import (
	"opencourse/common"
	"opencourse/database"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ConnectionString = "mongodb://localhost"

func getContext() *database.DbContext {

	context := &database.DbContext{}

	// Init default values
	context.Defaults()

	return context
}

func getAddCourseQuery() common.AddCourseQuery {
	addCourseQuery := common.AddCourseQuery{
		Name:       "The greatest golang",
		CategoryId: primitive.NewObjectID().Hex(),
		Tags:       []string{"Go"},
	}

	return addCourseQuery
}

func TestMain(m *testing.M) {
	context := getContext()

	err := context.Connect(ConnectionString)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	err = context.ClearCourses()

	if err != nil {
		panic(err)
	}

	m.Run()

	err = context.ClearCourses()

	if err != nil {
		panic(err)
	}
}

// TestAddCourse
func TestAddCourse(t *testing.T) {

	context := getContext()

	err := context.Connect(ConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	addCourseQuery := getAddCourseQuery()
	id, err := context.AddCourse(&addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(id)

}

// TestGetCourse
func TestGetCourse(t *testing.T) {

	context := getContext()

	err := context.Connect(ConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	addCourseQeury := getAddCourseQuery()

	id, err := context.AddCourse(&addCourseQeury)

	if err != nil {
		t.Fatal(err)
	}

	courseResult, err := context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(courseResult.Id)
}

// TestGetCourses
func TestGetCourses(t *testing.T) {
	context := getContext()

	err := context.Connect(ConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	for i := 0; i < 10; i++ {
		addCourseQeury := getAddCourseQuery()
		_, err := context.AddCourse(&addCourseQeury)
		if err != nil {
			t.Fatal(err)
		}
	}

	courses, err := context.GetCourses(getAddCourseQuery().CategoryId, 10, 0)

	if len(courses) < 10 {
		t.Error(err)
	}
}

// TestAddCourseTags
func TestAddCourseTags(t *testing.T) {
	context := getContext()

	err := context.Connect(ConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	addCourseQuery := getAddCourseQuery()
	id, err := context.AddCourse(&addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	err = context.AddCourseTags(id, []string{
		"C#",
		"C++",
		"Java",
		"Golang",
		"MongoDB",
		"PostgreSQL",
	})

	if err != nil {
		t.Fatal(err)
	}

	resultCourse, err := context.GetCourse(id)

	if len(resultCourse.Tags) != 7 {
		t.Error("course must contains 6 tags")
	}
}

// TestRemoveCourseTags
func TestRemoveCourseTags(t *testing.T) {
	context := getContext()

	err := context.Connect(ConnectionString)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	addCourseQuery := getAddCourseQuery()
	id, err := context.AddCourse(&addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	err = context.AddCourseTags(id, []string{
		"C#",
		"C++",
		"Java",
		"Golang",
		"MongoDB",
		"PostgreSQL",
	})

	if err != nil {
		t.Fatal(err)
	}

	err = context.RemoveCourseTags(id, []string{
		"C#",
		"Golang",
		"PostgreSQL",
	})

	resultCourse, err := context.GetCourse(id)

	if len(resultCourse.Tags) != 4 {
		t.Error("course must contains 4 tags")
	}

	for _, tag := range resultCourse.Tags {
		if tag == "C#" || tag == "Golang" || tag == "PostgreSQL" {
			t.Error("the selected tests are not deleted")
		}
	}
}
