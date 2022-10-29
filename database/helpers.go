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
ToAction map MgAction to Action
*/
func (mgAction *MgAction) ToAction() (*common.Action, error) {
	if mgAction == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToAction",
			},
			Model: "mgAction",
		}
	}

	var action common.Action
	action.ActionType = mgAction.ActionType
	action.UserId = mgAction.UserId.Hex()

	return &action, nil
}

/*
ToComment map MgComment to Comment
*/
func (mgComment *MgComment) ToComment() (*common.Comment, error) {
	if mgComment == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToComment",
			},
			Model: "mgComment",
		}
	}

	var comment common.Comment

	comment.Id = mgComment.Id.Hex()
	comment.UserId = mgComment.UserId.Hex()
	comment.Text = mgComment.Text
	if primitive.ObjectID.IsZero(mgComment.ParentId) == false {
		comment.ParentId = mgComment.ParentId.Hex()
	}

	if mgComment.Actions != nil {
		comment.Actions = []*common.Action{}

		for _, mgAction := range mgComment.Actions {
			action, err := mgAction.ToAction()

			if err == nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "database/mongodb/helpers.go",
						Method: "ToComment",
					},
					Msg: "can't convert MgAction to Action",
				}
			}

			comment.Actions = append(comment.Actions, action)
		}
	}

	return &comment, nil
}

/*
ToCourse map MgCourse to Course
*/
func (mgCourse *MgCourse) ToCourse() (*common.Course, error) {
	if mgCourse == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToCourse",
			},
			Model: "mgCourse",
		}
	}

	var course common.Course

	course.Id = mgCourse.Id.Hex()
	course.CategoryId = mgCourse.CategoryId.Hex()
	course.SubCategoryNumber = mgCourse.SubCategoryNumber
	course.Name = mgCourse.Name
	course.Lang = mgCourse.Lang
	course.Tags = mgCourse.Tags
	course.HeaderImg = mgCourse.HeaderImg
	course.Description = mgCourse.Description

	if mgCourse.Comments != nil {
		course.Comments = []*common.Comment{}

		for _, mgComment := range mgCourse.Comments {

			comment, err := mgComment.ToComment()

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

	if mgCourse.Actions != nil {
		course.Actions = []*common.Action{}

		for _, mgAction := range mgCourse.Actions {

			action, err := mgAction.ToAction()

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

	course.HeaderImg = mgCourse.HeaderImg

	course.Rating = mgCourse.Rating
	course.DateCreate = mgCourse.DateCreate.Time()
	course.DateUpdate = mgCourse.DateUpdate.Time()

	return &course, nil
}

/*
ToCategory map MgCategory to Category
*/
func (mgCategory *MgCategory) ToCategory() (*common.Category, error) {
	if mgCategory == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToCategory",
			},
			Model: "mgCategory",
		}
	}

	var category common.Category

	category.Id = mgCategory.Id.Hex()
	category.Name = mgCategory.Name
	category.Lang = mgCategory.Lang

	for _, mgSubCategory := range mgCategory.SubCategories {

		subCategory, err := mgSubCategory.ToSubCategory()

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
ToSubCategory map MgSubCategory to SubCategory
*/
func (mgSubCategory *MgSubCategory) ToSubCategory() (*common.SubCategory, error) {
	if mgSubCategory == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/helpers.go",
				Method: "ToSubCategory",
			},
			Model: "mgSubCategory",
		}
	}

	var subCategory common.SubCategory
	subCategory.Number = mgSubCategory.Number
	subCategory.Name = mgSubCategory.Name

	return &subCategory, nil

}
