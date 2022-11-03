package database

import (
	"opencourse/common"
	"opencourse/common/openerrors"
)

/*
ToCourse map DbCourse to Course
*/
func (dbCourse *DbCourse) ToCourse() (*common.Course, error) {
	if dbCourse == nil {
		return nil, openerrors.ModelNilOrEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToCourse",
			},
			Model: "dbCourse",
		}
	}

	var course common.Course

	course.Id = dbCourse.Id.Hex()
	course.CategoryId = dbCourse.CategoryId.Hex()
	course.Name = dbCourse.Name
	course.Tags = dbCourse.Tags
	course.Description = dbCourse.Description
	course.Rating = dbCourse.Rating
	course.Enabled = dbCourse.Enabled
	course.IconImg = dbCourse.IconImg
	course.HeaderImg = dbCourse.HeaderImg
	course.DateCreate = dbCourse.DateCreate.Time()
	course.DateUpdate = dbCourse.DateUpdate.Time()

	return &course, nil
}

/*
ToCategory map DbCategory to Category
*/
func (dbCategory *DbCategory) ToCategory() (*common.Category, error) {
	if dbCategory == nil {
		return nil, openerrors.ModelNilOrEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToCategory",
			},
			Model: "dbCategory",
		}
	}

	var category common.Category

	category.Id = dbCategory.Id.Hex()
	category.Name = dbCategory.Name
	category.Lang = dbCategory.Lang
	category.IconImg = dbCategory.IconImg
	category.HeaderImg = dbCategory.HeaderImg

	return &category, nil
}

/*
ToStage map DbStage to Stage
*/
func (dbStage *DbStage) ToStage() (*common.Stage, error) {
	if dbStage == nil {
		return nil, openerrors.ModelNilOrEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToStage",
			},
			Model: "dbStage",
		}
	}

	var stage common.Stage

	stage.Id = dbStage.Id.Hex()
	stage.Name = dbStage.Name
	stage.CourseId = dbStage.CourseId.Hex()
	stage.HeaderImg = dbStage.HeaderImg
	stage.OrderNumber = dbStage.OrderNumber

	if dbStage.Content != nil {
		stage.Content = &common.PostContent{}
		stage.Content.Body = dbStage.Content.Body
		stage.Content.MediaItems = dbStage.Content.MediaItems
	}

	return &stage, nil
}

/*
ToStagePreview map DbStage to StagePreview
*/
func (dbStage *DbStage) ToStagePreview() (*common.StagePreview, error) {
	if dbStage == nil {
		return nil, openerrors.ModelNilOrEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToStagePreview",
			},
			Model: "dbStage",
		}
	}

	var stage common.StagePreview

	stage.Id = dbStage.Id.Hex()
	stage.Name = dbStage.Name
	stage.CourseId = dbStage.CourseId.Hex()
	stage.HeaderImg = dbStage.HeaderImg
	stage.OrderNumber = dbStage.OrderNumber

	return &stage, nil
}

/*
ToTest map DbTest to Test
*/
func (dbTest *DbTest) ToTest() (*common.Test, error) {
	if dbTest == nil {
		return nil, openerrors.ModelNilOrEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToTest",
			},
			Model: "dbTest",
		}
	}

	var test common.Test

	test.Id = dbTest.Id.Hex()
	test.StageId = dbTest.StageId.Hex()
	test.TestType = dbTest.TestType
	test.LemmingsCount = dbTest.LemmingsCount
	test.OrderNumber = dbTest.OrderNumber

	if dbTest.OptionTest != nil {
		test.OptionTest = &common.OptionTest{}
		test.OptionTest.Question = dbTest.OptionTest.Question

		for _, dbOption := range dbTest.OptionTest.Options {
			test.OptionTest.Options =
				append(test.OptionTest.Options, &common.Option{Answer: dbOption.Answer, IsRight: dbOption.IsRight})
		}
	}

	if dbTest.RewriteTest != nil {
		test.RewriteTest = &common.RewriteTest{
			Question:    dbTest.RewriteTest.Question,
			RightAnswer: dbTest.RewriteTest.RightAnswer,
		}
	}

	return &test, nil
}

/*
ToTestPreview map DbTest to TestPreview
*/
func (dbTest *DbTest) ToTestPreview() (*common.TestPreview, error) {
	if dbTest == nil {
		return nil, openerrors.ModelNilOrEmptyErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToTestPreview",
			},
			Model: "dbTest",
		}
	}

	var test common.TestPreview

	test.Id = dbTest.Id.Hex()
	test.StageId = dbTest.StageId.Hex()
	test.TestType = dbTest.TestType
	test.LemmingsCount = dbTest.LemmingsCount
	test.OrderNumber = dbTest.OrderNumber

	return &test, nil
}
