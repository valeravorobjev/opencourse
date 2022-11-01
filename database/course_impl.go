package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common"
	"opencourse/common/openerrors"
)

// ClearCourses remove all data from course collection
func (ctx *DbContext) ClearCourses() error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	_, err := col.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ClearCourseCollection",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
GetCourse return course from db by id. Parameters:
id - course id;
*/
func (ctx *DbContext) GetCourse(id string) (*common.Course, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "GetCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	filter := bson.D{
		{"_id", objectId},
	}

	var dbCourse DbCourse
	err = col.FindOne(context.Background(), filter).Decode(&dbCourse)

	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "GetCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	course, err := dbCourse.ToCourse()

	if err != nil {
		return nil, openerrors.OpenDefaultErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "GetCourse",
			},
			Msg: err.Error(),
		}
	}

	return course, nil
}

/*
GetCourses return courses from db. Parameters:
take - how many records need to return;
skip - how many records deed to skip;

TODO: Now, courses range by rating and date. Need to upgrade rage system
*/
func (ctx *DbContext) GetCourses(take int64, skip int64) ([]*common.Course, error) {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	ops := options.Find().SetLimit(take).SetSkip(skip).
		SetSort(bson.D{{"rating", -1}, {"date_update", -1}})

	cursor, err := col.Find(context.Background(), bson.D{}, ops)

	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "GetCourses",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var dbCourses []*DbCourse

	err = cursor.All(context.Background(), &dbCourses)

	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "GetCourses",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var courses []*common.Course

	for _, dbCourse := range dbCourses {
		course, err := dbCourse.ToCourse()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "database/mongodb/course_impl.go",
					Method: "GetCourses",
				},
				Msg: err.Error(),
			}
		}

		courses = append(courses, course)
	}

	return courses, nil
}

/*
AddCourse add course to db. Parameters:
addCourseQuery - parameter for create new course;
*/
func (ctx *DbContext) AddCourse(addCourseQuery *common.AddCourseQuery) (string, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	var dbCourse DbCourse
	dbCourse.Name = addCourseQuery.Name

	objectCategoryId, err := primitive.ObjectIDFromHex(addCourseQuery.CategoryId)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	dbCourse.CategoryId = objectCategoryId
	dbCourse.Tags = addCourseQuery.Tags
	dbCourse.Rating = 0

	result, err := col.InsertOne(context.Background(), dbCourse)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
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
AddCourseTags - add tags to course. Parameters:
id - course id;
tags - tags;
*/
func (ctx *DbContext) AddCourseTags(id string, tags []string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	find := bson.D{
		{"_id", objectId},
	}

	update := bson.D{
		{"$push", bson.D{
			{
				"tags", bson.D{{
					"$each", tags,
				}},
			},
		}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
RemoveCourseTags - remove tags from course. Parameters:
id - course id;
tags - tags;
*/
func (ctx *DbContext) RemoveCourseTags(id string, tags []string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	find := bson.D{
		{"_id", objectId},
	}

	update := bson.D{
		{"$pullAll", bson.D{
			{
				"tags", tags,
			},
		}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
