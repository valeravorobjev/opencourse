package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"opencourse/common"
	"opencourse/common/openerrors"
)

// ClearTests remove all data from tests collection
func (ctx *DbContext) ClearTests() error {
	col := ctx.Client.Database(DbName).Collection(TestsCollection)

	_, err := col.DeleteMany(context.Background(), bson.D{{}})

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/stage_impl.go",
				Method: "ClearTests",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
GetTest return test. Parameters:
testId - test id;
*/
func (ctx *DbContext) GetTest(testId string) (*common.Test, error) {

	col := ctx.Client.Database(DbName).Collection(TestsCollection)

	objectTestId, err := primitive.ObjectIDFromHex(testId)

	if err != nil {
		return nil, openerrors.InvalidIdErr{
			Id:        testId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/test_impl.go",
					Method: "GetTest",
				},
				Msg: err.Error(),
			},
		}
	}

	filter := bson.D{{"_id", objectTestId}}

	var dbTest DbTest
	err = col.FindOne(context.Background(), filter).Decode(&dbTest)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "GetTest",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	test, err := dbTest.ToTest()

	if err != nil {
		return nil, openerrors.DefaultErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "GetTest",
			},
			Msg: err.Error(),
		}
	}

	return test, nil
}

/*
GetTests return tests. Parameters:
stageId - stage id;
take - how much records take;
skip - how much records skip;
*/
func (ctx *DbContext) GetTests(stageId string, take int64, skip int64) ([]*common.TestPreview, error) {
	col := ctx.Client.Database(DbName).Collection(TestsCollection)

	objectStageId, err := primitive.ObjectIDFromHex(stageId)

	if err != nil {
		return nil, openerrors.InvalidIdErr{
			Id:        stageId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/test_impl.go",
					Method: "GetTests",
				},
				Msg: err.Error(),
			},
		}
	}

	ops := options.Find().SetLimit(take).SetSkip(skip).
		SetSort(bson.D{{"order_number", 1}}).SetProjection(bson.D{{"option_test", -1}, {"rewrite_test", -1}})

	cursor, err := col.Find(context.Background(), bson.D{{"course_id", objectStageId}}, ops)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "GetTests",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var dbTests []*DbTest

	err = cursor.All(context.Background(), &dbTests)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "GetTests",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var tests []*common.TestPreview

	for _, dbTest := range dbTests {

		test, err := dbTest.ToTestPreview()

		if err != nil {
			return nil, openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/test_impl.go",
					Method: "GetTests",
				},
				Msg: err.Error(),
			}
		}

		tests = append(tests, test)
	}

	return tests, nil
}

/*
AddTest return tests. Parameters:
query - model for create test;
*/
func (ctx *DbContext) AddTest(query *common.AddTestQuery) (string, error) {
	col := ctx.Client.Database(DbName).Collection(TestsCollection)

	if query.OrderNumber < 0 {
		return "", openerrors.FieldEmptyErr{
			Field: "query.OrderNumber",
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "AddTest",
			},
		}
	}

	if len(query.TestType) < 2 {
		return "", openerrors.FieldEmptyErr{
			Field: "query.TestType",
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "AddTest",
			},
		}
	}

	if query.LemmingsCount < 1 {
		return "", openerrors.FieldEmptyErr{
			Field: "query.LemmingsCount",
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "AddTest",
			},
		}
	}

	objectStageId, err := primitive.ObjectIDFromHex(query.StageId)

	if err != nil {
		return "", openerrors.InvalidIdErr{
			Id:        query.StageId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/test_impl.go",
					Method: "AddTest",
				},
				Msg: err.Error(),
			},
		}
	}

	var dbTest DbTest

	dbTest.StageId = objectStageId
	dbTest.TestType = query.TestType
	dbTest.LemmingsCount = query.LemmingsCount
	dbTest.OrderNumber = query.OrderNumber

	if query.TestType == common.TestOption {

		if query.OptionTest == nil {
			return "", openerrors.FieldEmptyErr{
				Field: "query.OptionTest",
				BaseErr: openerrors.BaseErr{
					File:   "database/test_impl.go",
					Method: "AddTest",
				},
			}
		}

		dbTest.OptionTest = &DbOptionTest{}
		dbTest.OptionTest.Question = query.OptionTest.Question
		for _, option := range query.OptionTest.Options {
			dbTest.OptionTest.Options =
				append(dbTest.OptionTest.Options, &DbOption{Answer: option.Answer, IsRight: option.IsRight})
		}
	}

	if query.TestType == common.TestRewrite {

		if query.RewriteTest == nil {
			return "", openerrors.FieldEmptyErr{
				Field: "query.RewriteTest",
				BaseErr: openerrors.BaseErr{
					File:   "database/test_impl.go",
					Method: "AddTest",
				},
			}
		}

		dbTest.RewriteTest = &DbRewriteTest{}
		dbTest.RewriteTest.Question = query.RewriteTest.Question
		dbTest.RewriteTest.RightAnswer = query.RewriteTest.RightAnswer
	}

	result, err := col.InsertOne(context.Background(), dbTest)

	if err != nil {
		return "", openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "AddTest",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil

}

/*
DeleteTest delete test for stage. Parameters:
testId - test id;
*/
func (ctx *DbContext) DeleteTest(testId string) error {
	col := ctx.Client.Database(DbName).Collection(TestsCollection)

	ojbectTestId, err := primitive.ObjectIDFromHex(testId)

	if err != nil {
		return openerrors.InvalidIdErr{
			Id:        testId,
			Converter: "ObjectIDFromHex",
			Default: openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/test_impl.go",
					Method: "DeleteTest",
				},
				Msg: err.Error(),
			},
		}
	}

	_, err = col.DeleteOne(context.Background(), bson.D{{"_id", ojbectTestId}})

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/test_impl.go",
				Method: "DeleteTest",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
