package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common/openerrors"
)

/*
GetCourse return course from db by id
@cid - course id
*/
func (ctx *ContextMongoDb) GetCourse(cid string) (*Course, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

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

/*
GetCourses return courses from db
@take - how many records need to return
@skip - how many records deed to skip

TODO: Now, courses range by rating and date. Need to upgrade rage system
*/
func (ctx *ContextMongoDb) GetCourses(take int64, skip int64) ([]*Course, error) {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	ops := options.Find().SetLimit(take).SetSkip(skip).
		SetSort(bson.D{{"rating", -1}, {"date_update", -1}})

	cursor, err := col.Find(context.Background(), bson.D{}, ops)

	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "GetCourses",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var courses []*Course

	err = cursor.Decode(courses)

	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "GetCourses",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return courses, nil
}

/*
AddCourse add course to db
@course - entity for save to db
*/
func (ctx *ContextMongoDb) AddCourse(course *Course) (string, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	result, err := col.InsertOne(context.Background(), course)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), err
}

/*
AddCourseAuthors add authors to course
@id - course id
@aids - authors ids
*/
func (ctx *ContextMongoDb) AddCourseAuthors(id string, aids []string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	filter := bson.D{{"_id", id}}

	update := bson.D{
		{"$push", bson.E{Key: "author_ids", Value: aids}},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseAuthors",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
RemoveCourseAuthors add authors to course
@id - course id
@aids - authors ids
*/
func (ctx *ContextMongoDb) RemoveCourseAuthors(id string, aids []string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	filter := bson.D{{"_id", id}}

	update := bson.D{
		{"$pullAll", bson.E{Key: "author_ids", Value: aids}},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseAuthors",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
