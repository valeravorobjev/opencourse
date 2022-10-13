package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"opencourse/data-providers/mongodb"
	"testing"
	"time"
)

func getContext() *mongodb.ContextMongoDb {
	context := &mongodb.ContextMongoDb{}

	// Init default values
	context.Defaults()

	return context
}

func getCourse() mongodb.Course {
	dateCreate := primitive.Timestamp{T: uint32(time.Now().Unix())}
	course := mongodb.Course{
		Names: []*mongodb.GlobStr{
			{Lang: mongodb.LangEn, Text: "This is a test"},
			{Lang: mongodb.LangFr, Text: "C' est un test"},
		},
		Authors:     []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()},
		Category:    primitive.NewObjectID(),
		SubCategory: 12,
		BannerImg:   "",
		Tags: []*mongodb.GlobStr{
			{Text: "Test", Lang: mongodb.LangEn},
		},
		DateCreate: dateCreate,
		DateUpdate: dateCreate,
		Rating:     0,
		Actions:    []*mongodb.Action{},
		Comments:   []*mongodb.Comment{},
	}

	return course
}

// TestAddCourse
func TestAddCourse(t *testing.T) {

	context := getContext()

	err := context.Connect("mongodb://localhost")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	course := getCourse()
	id, err := context.AddCourse(&course)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(id)

}

// TestGetCourse
func TestGetCourse(t *testing.T) {

	context := getContext()

	err := context.Connect("mongodb://localhost")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	course := getCourse()

	id, err := context.AddCourse(&course)

	if err != nil {
		t.Fatal(err)
	}

	courseResult, err := context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(courseResult.Id.String())
}

// TestGetCourses
func TestGetCourses(t *testing.T) {
	context := getContext()

	err := context.Connect("mongodb://localhost")
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
		course := getCourse()
		_, err := context.AddCourse(&course)
		if err != nil {
			t.Fatal(err)
		}
	}

	courses, err := context.GetCourses(10, 0)

	if len(courses) < 10 {
		t.Error(err)
	}
}

// TestAddCourseAuthors
func TestAddCourseAuthors(t *testing.T) {
	context := getContext()

	err := context.Connect("mongodb://localhost")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	course := getCourse()
	id, err := context.AddCourse(&course)

	if err != nil {
		t.Fatal(err)
	}

	authorIds := []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()}

	err = context.AddCourseAuthors(id, authorIds)

	if err != nil {
		t.Fatal(err)
	}
}

// TestRemoveCourseAuthors
func TestRemoveCourseAuthors(t *testing.T) {
	context := getContext()

	err := context.Connect("mongodb://localhost")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Fatal(err)
		}
	}()

	course := getCourse()
	id, err := context.AddCourse(&course)

	if err != nil {
		t.Fatal(err)
	}

	authorIds := []string{course.Authors[0].Hex(), course.Authors[1].Hex()}

	err = context.RemoveCourseAuthors(id, authorIds)

	if err != nil {
		t.Fatal(err)
	}

	resultCourse, err := context.GetCourse(course.Id.Hex())

	if err != nil {
		t.Fatal(err)
	}

	if len(resultCourse.Authors) > 0 {
		t.Error("Authors field must be empty")
	}
}
