package database

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"opencourse/common"
	"opencourse/common/openerrors"
)

/*
This file contains methods for mapping common (Open) models to MongoDB models
*/

/*
ToAction map DbAction to Action
*/
func (dbAction *DbAction) ToAction() (*common.Action, error) {
	if dbAction == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToAction",
			},
			Model: "dbAction",
		}
	}

	var action common.Action
	action.ActionType = dbAction.ActionType
	action.UserId = dbAction.UserId.Hex()

	return &action, nil
}

/*
ToComment map DbComment to Comment
*/
func (dbComment *DbComment) ToComment() (*common.Comment, error) {
	if dbComment == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToComment",
			},
			Model: "dbComment",
		}
	}

	var comment common.Comment

	comment.Id = dbComment.Id.Hex()
	comment.UserId = dbComment.UserId.Hex()
	comment.Text = dbComment.Text
	if primitive.ObjectID.IsZero(dbComment.ParentId) == false {
		comment.ParentId = dbComment.ParentId.Hex()
	}

	if dbComment.Actions != nil {
		comment.Actions = []*common.Action{}

		for _, dbAction := range dbComment.Actions {
			action, err := dbAction.ToAction()

			if err == nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "database/mongodb/helpers.go",
						Method: "ToComment",
					},
					Msg: "can't convert DbAction to Action",
				}
			}

			comment.Actions = append(comment.Actions, action)
		}
	}

	return &comment, nil
}

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
	course.SubCategoryNumber = dbCourse.SubCategoryNumber
	course.Name = dbCourse.Name
	course.Lang = dbCourse.Lang
	course.Tags = dbCourse.Tags
	course.HeaderImg = dbCourse.HeaderImg
	course.Description = dbCourse.Description

	if dbCourse.Comments != nil {
		course.Comments = []*common.Comment{}

		for _, dbComment := range dbCourse.Comments {

			comment, err := dbComment.ToComment()

			if err != nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "database/mongodb/helpers.go",
						Method: "ToCourse",
					},
					Msg: err.Error(),
				}
			}

			course.Comments = append(course.Comments, comment)
		}

	}

	if dbCourse.Actions != nil {
		course.Actions = []*common.Action{}

		for _, dbAction := range dbCourse.Actions {

			action, err := dbAction.ToAction()

			if err != nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "database/mongodb/helpers.go",
						Method: "ToCourse",
					},
					Msg: err.Error(),
				}
			}

			course.Actions = append(course.Actions, action)
		}
	}

	course.HeaderImg = dbCourse.HeaderImg

	course.Rating = dbCourse.Rating
	course.DateCreate = dbCourse.DateCreate.Time()
	course.DateUpdate = dbCourse.DateUpdate.Time()

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

	for _, dbSubCategory := range dbCategory.SubCategories {

		subCategory, err := dbSubCategory.ToSubCategory()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "database/mongodb/helpers.go",
					Method: "ToCategory",
				},
				Msg: err.Error(),
			}
		}

		category.SubCategories = append(category.SubCategories, subCategory)
	}

	return &category, nil
}

/*
ToSubCategory map DbSubCategory to SubCategory
*/
func (dbSubCategory *DbSubCategory) ToSubCategory() (*common.SubCategory, error) {
	if dbSubCategory == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToSubCategory",
			},
			Model: "dbSubCategory",
		}
	}

	var subCategory common.SubCategory
	subCategory.Number = dbSubCategory.Number
	subCategory.Name = dbSubCategory.Name

	return &subCategory, nil

}
