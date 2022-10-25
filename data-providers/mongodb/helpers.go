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
ToOpenGlobStr map GlobStr to OpenGlobStr
*/
func (globStr *GlobStr) ToOpenGlobStr() (*common.OpenGlobStr, error) {
	if globStr == nil {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenGlobStr",
			},
			Model: "globStr",
		}
	}

	openGlobStr := common.OpenGlobStr{Lang: globStr.Lang, Text: globStr.Text}

	return &openGlobStr, nil
}

/*
ToOpenGlobStrs map GlobStr array to OpenGlobStr array
*/
func ToOpenGlobStrs(globStrs []*GlobStr) ([]*common.OpenGlobStr, error) {
	if globStrs == nil || len(globStrs) == 0 {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenGlobStrs",
			},
			Model: "globStrs",
		}
	}

	var openGlobStrs []*common.OpenGlobStr

	for _, globStr := range globStrs {

		openGlobStr, err := globStr.ToOpenGlobStr()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/helpers.go",
					Method: "ToOpenGlobStrs",
				},
				Msg: err.Error(),
			}
		}

		openGlobStrs = append(openGlobStrs, openGlobStr)
	}

	return openGlobStrs, nil
}

/*
ToGlobStrs map OpenGlobStr array to GlobStr array
*/
func ToGlobStrs(openGlobStrs []*common.OpenGlobStr) ([]*GlobStr, error) {
	if openGlobStrs == nil || len(openGlobStrs) == 0 {
		return nil, openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToGlobStrs",
			},
			Model: "openGlobStrs",
		}
	}

	var globStrs []*GlobStr

	for _, openGlobStr := range openGlobStrs {

		globStr := &GlobStr{}

		err := globStr.ToGlobStr(openGlobStr)

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/helpers.go",
					Method: "ToGlobStrs",
				},
				Msg: err.Error(),
			}
		}

		globStrs = append(globStrs, globStr)
	}

	return globStrs, nil
}

/*
ToGlobStr map OpenGlobStr to GlobStr
*/
func (globStr *GlobStr) ToGlobStr(openGlobStr *common.OpenGlobStr) error {
	if openGlobStr == nil {
		return openerrors.OpenModelNilOrEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToGlobStr",
			},
			Model: "openGlobStr",
		}
	}

	globStr.Lang = openGlobStr.Lang
	globStr.Text = openGlobStr.Text

	return nil
}

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
	openCourse.Names = []*common.OpenGlobStr{}

	if course.Names == nil || len(course.Names) == 0 {
		return nil, openerrors.OpenFieldEmptyErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/helpers.go",
				Method: "ToOpenCourse",
			},
			Field: "course.Names",
		}
	}

	for _, cn := range course.Names {
		openGlobStr, err := cn.ToOpenGlobStr()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/helpers.go",
					Method: "ToOpenCourse",
				},
				Msg: err.Error(),
			}
		}

		openCourse.Names = append(openCourse.Names, openGlobStr)
	}

	if course.Tags != nil {
		openCourse.Tags = []*common.OpenGlobStr{}

		for _, tag := range course.Tags {

			openGlobStr, err := tag.ToOpenGlobStr()

			if err != nil {
				return nil, openerrors.OpenDefaultErr{
					BaseErr: openerrors.OpenBaseErr{
						File:   "data-providers/mongodb/helpers.go",
						Method: "ToOpenCourse",
					},
					Msg: err.Error(),
				}
			}

			openCourse.Tags = append(openCourse.Tags, openGlobStr)
		}
	}

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

	if course.AuthorIds != nil {
		openCourse.AuthorIds = []string{}

		for _, id := range course.AuthorIds {
			openCourse.AuthorIds = append(openCourse.AuthorIds, id.Hex())
		}
	}

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

	for _, name := range category.Names {

		openName, err := name.ToOpenGlobStr()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/helpers.go",
					Method: "ToOpenCategory",
				},
				Msg: err.Error(),
			}
		}

		openCategory.Names = append(openCategory.Names, openName)
	}

	openCategory.Langs = category.Langs

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

	for _, name := range subCategory.Names {

		openName, err := name.ToOpenGlobStr()

		if err != nil {
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/helpers.go",
					Method: "ToOpenSubCategory",
				},
				Msg: err.Error(),
			}
		}

		openSubCategory.Names = append(openSubCategory.Names, openName)
	}

	return &openSubCategory, nil

}
