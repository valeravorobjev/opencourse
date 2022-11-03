package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common"
	"opencourse/common/openerrors"
	"time"
)

// ClearCourses remove all data from course collection
func (ctx *DbContext) ClearCourses() error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	_, err := col.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "ClearCourses",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
GetCourse return course from db by id. Parameters:
courseId - course id;
*/
func (ctx *DbContext) GetCourse(courseId string) (*common.Course, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectCourseId, err := primitive.ObjectIDFromHex(courseId)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "GetCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	filter := bson.D{
		{"_id", objectCourseId},
	}

	var dbCourse DbCourse
	err = col.FindOne(context.Background(), filter).Decode(&dbCourse)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "GetCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	course, err := dbCourse.ToCourse()

	if err != nil {
		return nil, openerrors.DefaultErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "GetCourse",
			},
			Msg: err.Error(),
		}
	}

	return course, nil
}

/*
GetCourses return courses from db. Parameters:
categoryId - category id. Optional, may be set empty string. Example: "".
take - how many records need to return;
skip - how many records need to skip;
*/
func (ctx *DbContext) GetCourses(categoryId string, take int64, skip int64) ([]*common.Course, error) {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	var filters []bson.D

	if len(categoryId) > 1 {
		objectCategoryId, err := primitive.ObjectIDFromHex(categoryId)

		if err != nil {
			return nil, openerrors.InvalidIdErr{
				Id:        categoryId,
				Converter: "ObjectIDFromHex",
				Default: openerrors.DefaultErr{
					BaseErr: openerrors.BaseErr{
						File:   "database/course_impl.go",
						Method: "GetCourses",
					},
					Msg: err.Error(),
				},
			}
		}

		filters = append(filters, bson.D{{"category_id", objectCategoryId}})
	}

	ops := options.Find().SetLimit(take).SetSkip(skip).
		SetSort(bson.D{{"rating", -1}, {"date_update", -1}})

	cursor, err := col.Find(context.Background(), filters, ops)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "GetCourses",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	var dbCourses []*DbCourse

	err = cursor.All(context.Background(), &dbCourses)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "GetCourses",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	var courses []*common.Course

	for _, dbCourse := range dbCourses {
		course, err := dbCourse.ToCourse()

		if err != nil {
			return nil, openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/course_impl.go",
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
		return "", openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "AddCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	dbCourse.CategoryId = objectCategoryId
	dbCourse.Tags = addCourseQuery.Tags
	dbCourse.Rating = 0

	dateNow := time.Now().UTC()
	dbCourse.DateCreate = primitive.NewDateTimeFromTime(dateNow)
	dbCourse.DateUpdate = primitive.NewDateTimeFromTime(dateNow)

	result, err := col.InsertOne(context.Background(), dbCourse)

	if err != nil {
		return "", openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "AddCourse",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
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
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
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
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
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
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
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
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.ConStr,
			DbErr:  err.Error(),
		}
	}

	return nil
}
