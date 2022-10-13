package integration

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

// TestAddCourse
func TestAddCourse(t *testing.T) {

	context := getContext()

	err := context.Connect("mongodb://localhost")
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = context.Disconnect()
		if err != nil {
			t.Error(err)
		}
	}()

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

	id, err := context.AddCourse(&course)

	if err != nil {
		t.Error(err)
	}

	t.Log(id)

}
