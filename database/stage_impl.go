package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common"
	"opencourse/common/openerrors"
)

// ClearStages remove all data from stages collection
func (ctx *DbContext) ClearStages() error {
	col := ctx.Client.Database(DbName).Collection(StageCollection)

	_, err := col.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "ClearStages",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
GetStage return stage by id. Parameters:
stageId - stage id;
*/
func (ctx *DbContext) GetStage(stageId string) (*common.Stage, error) {

	col := ctx.Client.Database(DbName).Collection(StageCollection)

	objectStageId, err := primitive.ObjectIDFromHex(stageId)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "GetStage",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	filter := bson.D{{"_id", objectStageId}}

	var dbStage DbStage
	err = col.FindOne(context.Background(), filter).Decode(&dbStage)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "GetStage",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	stage, err := dbStage.ToStage()

	if err != nil {
		return nil, openerrors.DefaultErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "GetStage",
			},
			Msg: err.Error(),
		}
	}

	return stage, nil
}

/*
GetStages return stages for course. Parameters:
courseId - course id;
take - how much records take;
skip - how much records skip;
*/
func (ctx *DbContext) GetStages(courseId string, take int64, skip int64) ([]*common.Stage, error) {
	col := ctx.Client.Database(DbName).Collection(StageCollection)

	objectCourseId, err := primitive.ObjectIDFromHex(courseId)

	if err != nil {
		return nil, openerrors.InvalidIdErr{
			Id:        courseId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/stage_impl.go",
					Method: "GetStages",
				},
				Msg: err.Error(),
			},
		}
	}

	ops := options.Find().SetLimit(take).SetSkip(skip).
		SetSort(bson.D{{"order_number", 1}})

	cursor, err := col.Find(context.Background(), bson.D{{"course_id", objectCourseId}}, ops)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "GetStages",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var dbStages []*DbStage

	err = cursor.All(context.Background(), &dbStages)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "GetStages",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var stages []*common.Stage

	for _, dbStage := range dbStages {

		stage, err := dbStage.ToStage()

		if err != nil {
			return nil, openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/stage_impl.go",
					Method: "GetStages",
				},
				Msg: err.Error(),
			}
		}

		stages = append(stages, stage)
	}

	return stages, nil
}

/*
AddStage add stage for course. Parameters:
addStageQuery - model for create new stage;
*/
func (ctx *DbContext) AddStage(addStageQuery *common.AddStageQuery) (string, error) {
	col := ctx.Client.Database(DbName).Collection(StageCollection)

	// if name < 2 letters length than name is empty
	if len(addStageQuery.Name) < 2 {
		return "", openerrors.FieldEmptyErr{
			Field: "addStageQuery.Name",
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "AddStage",
			},
		}
	}

	// if header image path length < 5 chars, field is empty
	if len(addStageQuery.HeaderImg) < 5 {
		return "", openerrors.FieldEmptyErr{
			Field: "addStageQuery.HeaderImg",
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "AddStage",
			},
		}
	}

	// if order number < 0, field is empty
	if addStageQuery.OrderNumber < 0 {
		return "", openerrors.FieldEmptyErr{
			Field: "addStageQuery.OrderNumber",
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "AddStage",
			},
		}
	}

	// if content is null
	if addStageQuery.Content == nil {
		return "", openerrors.FieldEmptyErr{
			Field: "addStageQuery.Content",
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "AddStage",
			},
		}
	}

	var dbStage DbStage

	dbStage.Name = addStageQuery.Name
	dbStage.HeaderImg = addStageQuery.HeaderImg
	dbStage.OrderNumber = addStageQuery.OrderNumber

	objectCourseId, err := primitive.ObjectIDFromHex(addStageQuery.CourseId)

	if err != nil {
		return "", openerrors.InvalidIdErr{
			Id:        addStageQuery.CourseId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/stage_impl.go",
					Method: "AddStage",
				},
				Msg: err.Error(),
			},
		}
	}

	dbStage.CourseId = objectCourseId
	dbStage.Content = &DbPostContent{}
	dbStage.Content.Body = addStageQuery.Content.Body
	dbStage.Content.MediaItems = addStageQuery.Content.MediaItems

	result, err := col.InsertOne(context.Background(), dbStage)

	if err != nil {
		return "", openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "AddStage",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	newId := result.InsertedID.(primitive.ObjectID)

	return newId.Hex(), nil
}
