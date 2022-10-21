package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"opencourse/common"
	"opencourse/data-providers/mongodb"
	"testing"
)

const MongoDbHost = "mongodb://localhost"

func getContext() *mongodb.ContextMongoDb {
	context := &mongodb.ContextMongoDb{}

	// Init default values
	context.Defaults()

	return context
}

func getAddCourseQuery() common.OpenAddCourseQuery {
	addCourseQuery := common.OpenAddCourseQuery{
		Names: []*common.OpenGlobStr{
			{Lang: common.OpenLangEn, Text: "This is a test"},
			{Lang: common.OpenLangFr, Text: "C' est un test"},
		},
		Langs:             []string{common.OpenLangEn, common.OpenLangFr},
		CategoryId:        primitive.NewObjectID().Hex(),
		SubCategoryNumber: 12,
		Tags: []*common.OpenGlobStr{
			{Text: "Test", Lang: common.OpenLangEn},
		},
		HeaderImg: "",
	}

	return addCourseQuery
}

func TestMain(m *testing.M) {
	context := getContext()

	err := context.Connect(MongoDbHost)
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

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

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

	addCourseQeury := getAddCourseQuery()

	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQeury)

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

	err := context.Connect(MongoDbHost)
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
		_, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQeury)
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

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

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

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	authorsIds := []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()}

	err = context.AddCourseAuthors(id, authorsIds)

	if err != nil {
		t.Fatal(err)
	}

	err = context.RemoveCourseAuthors(id, authorsIds)

	if err != nil {
		t.Fatal(err)
	}

	resultCourse, err := context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	if len(resultCourse.AuthorIds) != 1 {
		t.Error("AuthorIds must have one author")
	}
}

// TestAddCourseAction
func TestAddCourseAction(t *testing.T) {
	context := getContext()

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	userId := primitive.NewObjectID().Hex()
	err = context.AddCourseAction(id, userId, common.OpenActionLike)

	if err != nil {
		t.Fatal(err)
	}

	err = context.AddCourseAction(id, userId, common.OpenActionDislike)

	if err != nil {
		t.Fatal(err)
	}

	err = context.AddCourseAction(id, userId, common.OpenActionLike)

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

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	userId := primitive.NewObjectID().Hex()
	err = context.AddCourseAction(id, userId, common.OpenActionLike)

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

	err = context.RemoveCourseAction(id, userId)

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

// TestChangeCourseAction
func TestChangeCourseAction(t *testing.T) {
	context := getContext()

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	userId := primitive.NewObjectID().Hex()
	err = context.AddCourseAction(id, userId, common.OpenActionLike)

	if err != nil {
		t.Fatal(err)
	}

	err = context.ChangeCourseAction(id, userId, common.OpenActionDislike)

	if err != nil {
		t.Fatal(err)
	}

	courseResult, err := context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	if len(courseResult.Actions) != 1 || courseResult.Actions[0].ActionType != common.OpenActionDislike {
		t.Error("Course must = 1 and ActionType = ActionDislike")
	}
}

// TestAddCourseComment
func TestAddCourseComment(t *testing.T) {
	context := getContext()

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	userId := primitive.NewObjectID().Hex()

	_, err = context.AddCourseComment(id, userId, "My comment")

	if err != nil {
		t.Fatal(err)
	}

}

// TestReplyCourseComment
func TestReplyCourseComment(t *testing.T) {
	context := getContext()

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	userId := primitive.NewObjectID().Hex()

	commentId, err := context.AddCourseComment(id, userId, "My comment")

	if err != nil {
		t.Fatal(err)
	}

	err = context.ReplyCourseComment(id, userId, commentId, "reply reply reply")

	if err != nil {
		t.Fatal(err)
	}

	resultCourse, err := context.GetCourse(id)

	if err != nil {
		t.Fatal(err)
	}

	if len(resultCourse.Comments) != 2 {
		t.Error("course must contains 2 comments")
	}

	if resultCourse.Comments[0].Id != resultCourse.Comments[1].ParentId {
		t.Error("reply comment parentId != base comment id")
	}

}

// TestRemoveCourseComment
func TestRemoveCourseComment(t *testing.T) {
	context := getContext()

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	userId := primitive.NewObjectID().Hex()

	commentId, err := context.AddCourseComment(id, userId, "My comment")

	if err != nil {
		t.Fatal(err)
	}

	commentId, err = context.AddCourseComment(id, userId, "My comment")

	if err != nil {
		t.Fatal(err)
	}

	err = context.RemoveCourseComment(id, commentId)

	if err != nil {
		t.Fatal(err)
	}

}

// TestAddCourseTags
func TestAddCourseTags(t *testing.T) {
	context := getContext()

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	err = context.AddCourseTags(id, []*common.OpenGlobStr{
		{Lang: common.OpenLangEn, Text: "C#"},
		{Lang: common.OpenLangEn, Text: "C++"},
		{Lang: common.OpenLangEn, Text: "Java"},
		{Lang: common.OpenLangEn, Text: "Golang"},
		{Lang: common.OpenLangEn, Text: "MongoDB"},
		{Lang: common.OpenLangEn, Text: "PostgreSQL"},
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

	err := context.Connect(MongoDbHost)
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
	id, err := context.AddCourse(primitive.NewObjectID().Hex(), &addCourseQuery)

	if err != nil {
		t.Fatal(err)
	}

	err = context.AddCourseTags(id, []*common.OpenGlobStr{
		{Lang: common.OpenLangEn, Text: "C#"},
		{Lang: common.OpenLangEn, Text: "C++"},
		{Lang: common.OpenLangEn, Text: "Java"},
		{Lang: common.OpenLangEn, Text: "Golang"},
		{Lang: common.OpenLangEn, Text: "MongoDB"},
		{Lang: common.OpenLangEn, Text: "PostgreSQL"},
	})

	if err != nil {
		t.Fatal(err)
	}

	err = context.RemoveCourseTags(id, []*common.OpenGlobStr{
		{Lang: common.OpenLangEn, Text: "C#"},
		{Lang: common.OpenLangEn, Text: "Golang"},
		{Lang: common.OpenLangEn, Text: "PostgreSQL"},
	})

	resultCourse, err := context.GetCourse(id)

	if len(resultCourse.Tags) != 4 {
		t.Error("course must contains 4 tags")
	}

	for _, tag := range resultCourse.Tags {
		if tag.Lang == common.OpenLangEn && (tag.Text == "C#" || tag.Text == "Golang" || tag.Text == "PostgreSQL") {
			t.Error("the selected tests are not deleted")
		}
	}
}
