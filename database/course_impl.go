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
GetCourse return course from db by id
@id - course id
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
GetCourses return courses from db
@take - how many records need to return
@skip - how many records deed to skip

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
AddCourse add course to db
userId - user id how create new course. He also set to course authors
@course - entity for save to db
*/
func (ctx *DbContext) AddCourse(addCourseQuery *common.AddCourseQuery) (string, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	var dbCourse DbCourse
	dbCourse.Name = addCourseQuery.Name
	dbCourse.Lang = addCourseQuery.Lang

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
	dbCourse.SubCategoryNumber = addCourseQuery.SubCategoryNumber
	dbCourse.HeaderImg = addCourseQuery.HeaderImg
	dbCourse.Tags = addCourseQuery.Tags

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
	dbCourse.Rating = 0

	nowUtcTime := time.Now().UTC()

	dbCourse.DateCreate = primitive.NewDateTimeFromTime(nowUtcTime)
	dbCourse.DateUpdate = primitive.NewDateTimeFromTime(nowUtcTime)

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
AddCourseAction set action for course
@id - course id
@userId - user id
@actionType - type of action
*/
func (ctx *DbContext) AddCourseAction(id string, userId string, actionType string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	objectUserId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	action := DbAction{UserId: objectUserId, ActionType: actionType}

	filter := bson.D{
		{"_id", objectId},
		{"actions.user_id", bson.D{{"$ne", objectUserId}}},
	}

	update := bson.D{
		{"$push", bson.D{{"actions", action}}},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
RemoveCourseAction remove action from course
@id - course id
@userId - user id
*/
func (ctx *DbContext) RemoveCourseAction(id string, userId string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "RemoveCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	objectUserId, err := primitive.ObjectIDFromHex(userId)

	filter := bson.D{{"_id", objectId}}

	update := bson.D{
		{"$pull",
			bson.D{
				{"actions", bson.D{{"user_id", objectUserId}}},
			},
		},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "RemoveCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
ChangeCourseAction change course action
@id - course id
@userId - user id
@actionType - type of action
*/
func (ctx *DbContext) ChangeCourseAction(id string, userId string, actionType string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ChangeCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	objectUserId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ChangeCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	filter := bson.D{{"_id", objectId}, {"actions.user_id", objectUserId}}

	update := bson.D{{"$set", bson.D{{"actions.$.action_type", actionType}}}}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ChangeCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil

}

/*
AddCourseComment add comment to course or for another course comment
@id - course id
@userId - user id
@text - comment's text
*/
func (ctx *DbContext) AddCourseComment(id string, userId string, text string) (string, error) {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectUserId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	dbComment := &DbComment{Id: primitive.NewObjectID(), UserId: objectUserId, Text: text}

	find := bson.D{
		{"_id", objectId},
	}

	update := bson.D{
		{"$push", bson.D{{"comments", dbComment}}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "AddCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return dbComment.Id.Hex(), nil
}

/*
ReplyCourseComment reply course comment
@id - course id
@userId - user id
@commentId - comment id
@text - comment's text
*/
func (ctx *DbContext) ReplyCourseComment(id string, userId string, commentId string, text string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectUserId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ReplyCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ReplyCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	objectCommentId, err := primitive.ObjectIDFromHex(commentId)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ReplyCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	dbComment := &DbComment{Id: primitive.NewObjectID(), UserId: objectUserId, Text: text, ParentId: objectCommentId}

	find := bson.D{
		{"_id", objectId},
		{"comments.id", objectCommentId},
	}

	update := bson.D{
		{"$push", bson.D{
			{"comments", dbComment},
		}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "ReplyCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
RemoveCourseComment remove comment from course
@id - course id
@commentId - comment id
*/
func (ctx *DbContext) RemoveCourseComment(id string, commentId string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectCommentId, err := primitive.ObjectIDFromHex(commentId)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "RemoveCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "RemoveCourseComment",
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
		{"$pull",
			bson.D{
				{"comments", bson.D{
					{"$or", bson.A{
						bson.D{
							{"id", objectCommentId},
						},
						bson.D{
							{"parent_id", objectCommentId},
						},
					}}},
				}},
		},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/course_impl.go",
				Method: "RemoveCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
AddCourseTags - add tags to course
@id - course id
@tags - tags
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
RemoveCourseTags - remove tags from course
@id - course id
@tags - tags
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