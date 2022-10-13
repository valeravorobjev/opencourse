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

func TestMain(m *testing.M) {
	context := getContext()

	err := context.Connect("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	err = context.ClearCourseCollection()

	if err != nil {
		panic(err)
	}

	m.Run()

	//err = context.ClearCourseCollection()
	//
	//if err != nil {
	//	panic(err)
	//}
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

// TestAddCourseAction
func TestAddCourseAction(t *testing.T) {
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

	userId := primitive.NewObjectID()
	action := mongodb.Action{ActionType: mongodb.ActionLike, UserId: userId}

	err = context.AddCourseAction(id, &action)

	if err != nil {
		t.Fatal(err)
	}

	err = context.AddCourseAction(id, &action)

	if err != nil {
		t.Fatal(err)
	}

	action.ActionType = mongodb.ActionDislike
	err = context.AddCourseAction(id, &action)

	if err != nil {
		t.Fatal(err)
	}

	courseResult, err := context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	if len(courseResult.Actions) != 1 {
		t.Error("Course contains more than one action")
	}

}

// TestRemoveCourseAction
func TestRemoveCourseAction(t *testing.T) {
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

	userId := primitive.NewObjectID()
	err = context.AddCourseAction(id, &mongodb.Action{UserId: userId, ActionType: mongodb.ActionLike})

	if err != nil {
		t.Fatal(err)
	}

	courseResult, err := context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	if len(courseResult.Actions) != 1 {
		t.Error("Course must contains one action")
	}

	err = context.RemoveCourseAction(id, userId.Hex())

	if err != nil {
		t.Fatal(err)
	}

	courseResult, err = context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	if len(courseResult.Actions) > 0 {
		t.Error("Course must be empty")
	}
}
