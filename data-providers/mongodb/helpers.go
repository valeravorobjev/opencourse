package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"opencourse/common"
	"opencourse/common/openerrors"
	"time"
)

/*
This file contains methods for mapping common (Open) models to MongoDB models
*/

/*
ToOpenAction map Action to OpenAction
*/
func (action *Action) ToOpenAction() (*common.OpenAction, error) {
	if action == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenAction",
			},
			Model: "action",
		}
	}

	var openAction common.OpenAction
	openAction.ActionType = action.ActionType
	openAction.UserId = action.UserId.Hex()

	return &openAction, nil
}

/*
ToOpenComment map Comment to OpenComment
*/
func (comment *Comment) ToOpenComment() (*common.OpenComment, error) {
	if comment == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenComment",
			},
			Model: "comment",
		}
	}

	var openComment common.OpenComment
	openComment.Id = comment.Id.Hex()
	openComment.UserId = comment.UserId.Hex()
	openComment.Text = comment.Text
	if primitive.ObjectID.IsZero(comment.ParentId) != false {
		openComment.ParentId = comment.ParentId.Hex()
	}

	if comment.Actions != nil {
		openComment.Actions = []*common.OpenAction{}

		for _, action := range comment.Actions {
			openAction, err := action.ToOpenAction()

			if err == nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "data-providers/mongodb/helpers.go",
						Method: "ToOpenComment",
					},
					Msg: "can't convert Action to OpenAction",
				}
			}

			openComment.Actions = append(openComment.Actions, openAction)
		}
	}

	return &openComment, nil
}

/*
ToOpenCourse map Course to OpenCourse
*/
func (course *Course) ToOpenCourse() (*common.OpenCourse, error) {
	if course == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenCourse",
			},
			Model: "course",
		}
	}

	var openCourse common.OpenCourse

	openCourse.Id = course.Id.Hex()
	openCourse.CategoryId = course.CategoryId.Hex()
	openCourse.SubCategoryNumber = course.SubCategoryNumber
	openCourse.Name = course.Name
	openCourse.Lang = course.Lang
	openCourse.Tags = course.Tags

	if course.Comments != nil {
		openCourse.Comments = []*common.OpenComment{}

		for _, comment := range course.Comments {

			openComment, err := comment.ToOpenComment()

			if err != nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "data-providers/mongodb/helpers.go",
						Method: "ToOpenCourse",
					},
					Msg: err.Error(),
				}
			}

			openCourse.Comments = append(openCourse.Comments, openComment)
		}

	}

	if course.Actions != nil {
		openCourse.Actions = []*common.OpenAction{}

		for _, action := range course.Actions {

			openAction, err := action.ToOpenAction()

			if err != nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "data-providers/mongodb/helpers.go",
						Method: "ToOpenCourse",
					},
					Msg: err.Error(),
				}
			}

			openCourse.Actions = append(openCourse.Actions, openAction)
		}
	}

	openCourse.HeaderImg = course.HeaderImg

	openCourse.Rating = course.Rating
	openCourse.DateCreate = time.Unix(int64(course.DateCreate.T), 0)
	openCourse.DateUpdate = time.Unix(int64(course.DateUpdate.T), 0)

	return &openCourse, nil
}

/*
ToOpenCategory map Category to OpenCategory
*/
func (category *Category) ToOpenCategory() (*common.OpenCategory, error) {
	if category == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenCategory",
			},
			Model: "category",
		}
	}

	var openCategory common.OpenCategory

	openCategory.Id = category.Id.Hex()
	openCategory.Name = category.Name
	openCategory.Lang = category.Lang

	for _, subCategory := range category.SubCategories {

		openSubCategory, err := subCategory.ToOpenSubCategory()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/helpers.go",
					Method: "ToOpenCategory",
				},
				Msg: err.Error(),
			}
		}

		openCategory.SubCategories = append(openCategory.SubCategories, openSubCategory)
	}

	return &openCategory, nil
}

/*
ToOpenSubCategory map SubCategory to OpenSubCategory
*/
func (subCategory *SubCategory) ToOpenSubCategory() (*common.OpenSubCategory, error) {
	if subCategory == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenSubCategory",
			},
			Model: "subCategory",
		}
	}

	var openSubCategory common.OpenSubCategory
	openSubCategory.Number = subCategory.Number
	openSubCategory.Name = subCategory.Name

	return &openSubCategory, nil

}
