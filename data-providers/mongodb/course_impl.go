package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common/openerrors"
)

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
@action - user action
*/
func (ctx *ContextMongoDb) AddCourseAction(id string, action *Action) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "SetCourseAction",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	filter := bson.D{{"_id", objectId},
		{"actions.user_id", bson.D{{"$ne", action.UserId}}}}

	update := bson.D{
		{"$push", bson.D{{"actions", action}}},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "SetCourseAction",
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

	filter := bson.D{{"_id", id}}

	//TODO:: check work or not
	update := bson.D{
		{"$pull",
			bson.D{
				{"actions", bson.E{Key: "user_id", Value: userId}},
			},
		},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/course_impl.go",
				Method: "UnsetCourseAction",
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

	comment := &Comment{Id: primitive.NewObjectID(), UserId: objectUserId, Text: text}

	find := bson.D{
		{"_id", id},
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
RemoveCourseComment remove comment from course
@id - course id
@commentId - comment id
*/
func (ctx *ContextMongoDb) RemoveCourseComment(id string, commentId string) error {
	col := ctx.Client.Database(DbName).Collection(CourseCollection)

	find := bson.D{
		{"_id", id}, {"comments.id", commentId},
	}

	update := bson.D{
		{"$pop", bson.D{{"comments", 1}}},
	}

	_, err := col.UpdateOne(context.Background(), find, update)

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
