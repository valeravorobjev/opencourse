package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"opencourse/common/openerrors"
)

/*
GetCourse return course from db by id
@cid - course id
*/
func (ctx *ContextMongoDb) GetCourse(cid string) (*Course, error) {

	col := ctx.client.Database(DbName).Collection(CourseCollection)

	filter := bson.D{
		{"_id", cid},
	}

	var course Course
	err := col.FindOne(context.Background(), filter).Decode(&course)

	if err != nil && err != mongo.ErrNoDocuments {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "GetCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return &course, nil
}
