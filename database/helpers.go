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
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
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

	return &course, nil
}

/*
ToCategory map DbCategory to Category
*/
func (dbCategory *DbCategory) ToCategory() (*common.Category, error) {
	if dbCategory == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
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

	return &category, nil
}
