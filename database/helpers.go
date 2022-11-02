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
