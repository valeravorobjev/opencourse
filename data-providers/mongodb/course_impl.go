package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common"
	"opencourse/common/openerrors"
	"time"
)

// ClearCourseCollection remove all data from course collection
func (ctx *ContextMongoDb) ClearCourseCollection() error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	_, err := col.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
func (ctx *ContextMongoDb) GetCourse(id string) (*Course, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
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

	filter := bson.D{
		{"_id", objectId},
	}

	var course Course
	err = col.FindOne(context.Background(), filter).Decode(&course)

	if err != nil {
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
func (ctx *ContextMongoDb) GetCourses(take int64, skip int64) ([]*common.OpenCourse, error) {
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

	err = cursor.All(context.Background(), &courses)

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

	var openCourses []*common.OpenCourse

	for _, course := range courses {
		openCourse, err := course.ToOpenCourse()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/course_impl.go",
					Method: "GetCourses",
				},
				Msg: err.Error(),
			}
		}

		openCourses = append(openCourses, openCourse)
	}

	return openCourses, nil
}

/*
AddCourse add course to db
userId - user id how create new course. He also set to course authors
@course - entity for save to db
*/
func (ctx *ContextMongoDb) AddCourse(userId string, addCourseQuery *common.OpenAddCourseQuery) (string, error) {

	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	var course Course

	for _, name := range addCourseQuery.Names {
		courseName := &GlobStr{}

		err := courseName.ToGlobStr(name)

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

		course.Names = append(course.Names, courseName)
	}

	course.Langs = addCourseQuery.Langs

	categoryId, err := primitive.ObjectIDFromHex(addCourseQuery.CategoryId)

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

	course.CategoryId = categoryId
	course.SubCategoryNumber = addCourseQuery.SubCategoryNumber
	course.HeaderImg = addCourseQuery.HeaderImg

	if addCourseQuery.Tags != nil && len(addCourseQuery.Tags) > 0 {
		course.Tags = []*GlobStr{}

		for _, openTag := range addCourseQuery.Tags {
			tag := &GlobStr{}

			err = tag.ToGlobStr(openTag)

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

			course.Tags = append(course.Tags, tag)
		}
	}

	authorId, err := primitive.ObjectIDFromHex(userId)

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

	course.AuthorIds = []primitive.ObjectID{authorId}
	course.Rating = 0

	nowUtcTime := time.Now().Unix()

	course.DateCreate = primitive.Timestamp{T: uint32(nowUtcTime)}
	course.DateUpdate = primitive.Timestamp{T: uint32(nowUtcTime)}

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
	course.Id = id

	return id.Hex(), err
}

/*
AddCourseAuthors add authors to course
@id - course id
@aids - authors ids
*/
func (ctx *ContextMongoDb) AddCourseAuthors(id string, aids []string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

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

	var objectAids []primitive.ObjectID
	for _, aid := range aids {
		objectAid, err := primitive.ObjectIDFromHex(aid)

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

		objectAids = append(objectAids, objectAid)
	}

	filter := bson.D{{"_id", objectId}}

	update := bson.D{
		{"$push", bson.D{{"author_ids", bson.D{{"$each", objectAids}}}}},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

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
RemoveCourseAuthors delete authors from course
@id - course id
@aids - authors ids
*/
func (ctx *ContextMongoDb) RemoveCourseAuthors(id string, aids []string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "RemoveCourseAuthors",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var objectAids []primitive.ObjectID
	for _, aid := range aids {
		objectAid, err := primitive.ObjectIDFromHex(aid)

		if err != nil {
			return openerrors.OpenDbErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/course_impl.go",
					Method: "RemoveCourseAuthors",
				},
				DbName: ctx.DbName,
				ConStr: ctx.Uri,
				DbErr:  err.Error(),
			}
		}

		objectAids = append(objectAids, objectAid)
	}

	filter := bson.D{{"_id", objectId}}

	update := bson.D{
		{"$pullAll", bson.D{{"author_ids", objectAids}}},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "RemoveCourseAuthors",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
AddCourseAction set action for course
@id - course id
@userId - user id
@actionType - type of action
*/
func (ctx *ContextMongoDb) AddCourseAction(id string, userId string, actionType string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	action := Action{UserId: objectUserId, ActionType: actionType}

	filter := bson.D{{"_id", objectId},
		{"actions.user_id", bson.D{{"$ne", objectUserId}}}}

	update := bson.D{
		{"$push", bson.D{{"actions", action}}},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
func (ctx *ContextMongoDb) RemoveCourseAction(id string, userId string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
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
func (ctx *ContextMongoDb) ChangeCourseAction(id string, userId string, actionType string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
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
func (ctx *ContextMongoDb) AddCourseComment(id string, userId string, text string) (string, error) {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectUserId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	comment := &Comment{Id: primitive.NewObjectID(), UserId: objectUserId, Text: text}

	find := bson.D{
		{"_id", objectId},
	}

	update := bson.D{
		{"$push", bson.D{{"comments", comment}}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return comment.Id.Hex(), nil
}

/*
ReplyCourseComment reply course comment
@id - course id
@userId - user id
@commentId - comment id
@text - comment's text
*/
func (ctx *ContextMongoDb) ReplyCourseComment(id string, userId string, commentId string, text string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectUserId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
				Method: "ReplyCourseComment",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	comment := &Comment{Id: primitive.NewObjectID(), UserId: objectUserId, Text: text, ParentId: objectCommentId}

	find := bson.D{
		{"_id", objectId},
		{"comments.id", objectCommentId},
	}

	update := bson.D{
		{"$push", bson.D{
			{"comments", comment},
		}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
func (ctx *ContextMongoDb) RemoveCourseComment(id string, commentId string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectCommentId, err := primitive.ObjectIDFromHex(commentId)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
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
				File:   "data-providers/mongodb/course_impl.go",
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
func (ctx *ContextMongoDb) AddCourseTags(id string, tags []*common.OpenGlobStr) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var mongoTags []*GlobStr

	for _, tag := range tags {

		mongoTag := &GlobStr{}
		err = mongoTag.ToGlobStr(tag)

		if err != nil {
			return openerrors.OpenDbErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/course_impl.go",
					Method: "AddCourseTags",
				},
				DbName: ctx.DbName,
				ConStr: ctx.Uri,
				DbErr:  err.Error(),
			}
		}

		mongoTags = append(mongoTags, mongoTag)
	}

	find := bson.D{
		{"_id", objectId},
	}

	update := bson.D{
		{"$push", bson.D{
			{
				"tags", bson.D{{
					"$each", mongoTags,
				}},
			},
		}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
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
func (ctx *ContextMongoDb) RemoveCourseTags(id string, tags []*common.OpenGlobStr) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var mongoTags []*GlobStr

	for _, tag := range tags {

		mongoTag := &GlobStr{}
		err = mongoTag.ToGlobStr(tag)

		if err != nil {
			return openerrors.OpenDbErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/course_impl.go",
					Method: "AddCourseTags",
				},
				DbName: ctx.DbName,
				ConStr: ctx.Uri,
				DbErr:  err.Error(),
			}
		}

		mongoTags = append(mongoTags, mongoTag)
	}

	find := bson.D{
		{"_id", objectId},
	}

	update := bson.D{
		{"$pullAll", bson.D{
			{
				"tags", mongoTags,
			},
		}},
	}

	_, err = col.UpdateOne(context.Background(), find, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "AddCourseTags",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
