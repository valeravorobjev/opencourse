package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"opencourse/common"
	"opencourse/common/openerrors"
)

// GetCategories return all categories from db
func (ctx *DbContext) GetCategories(langs []string) ([]*common.Category, error) {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	find := bson.D{
		{"langs", bson.D{
			{
				"$in", langs,
			},
		}},
	}

	cursor, err := col.Find(context.Background(), find)

	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/category_impl.go",
				Method: "GetCategories",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var dbCategories []*DbCategory

	err = cursor.Decode(dbCategories)
	if err != nil {
		return nil, openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/category_impl.go",
				Method: "GetCategories",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var categories []*common.Category

	for _, dbCategory := range dbCategories {
		category, err := dbCategory.ToCategory()

		if err != nil {
			return nil, openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/category_impl.go",
					Method: "GetCategories",
				},
				Msg: err.Error(),
			}
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// AddCategory method for add new category and sub categories to database
func (ctx *DbContext) AddCategory(addCategoryQuery *common.AddCategoryQuery) (string, error) {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	category := DbCategory{}

	category.Name = addCategoryQuery.Name
	category.Lang = addCategoryQuery.Lang

	result, err := col.InsertOne(context.Background(), category)

	if err != nil {
		return "", openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/category_impl.go",
				Method: "AddCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	id := result.InsertedID.(primitive.ObjectID)

	return id.Hex(), err
}

/*
UpdateCategory update category.Parameters:
categoryId - category id;
name - category name;
lang - category language;
*/
func (ctx *DbContext) UpdateCategory(categoryId string, name string, lang string) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	objectCategoryId, err := primitive.ObjectIDFromHex(categoryId)

	if err != nil {
		if err != nil {
			return openerrors.DefaultErr{
				BaseErr: openerrors.BaseErr{
					File:   "database/category_impl.go",
					Method: "UpdateCategory",
				},
				Msg: err.Error(),
			}
		}
	}

	filter := bson.D{
		{"_id", objectCategoryId},
	}

	update := bson.D{
		{
			"$set", bson.D{{"name", name}, {"lang", lang}},
		},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.DbErr{
			BaseErr: openerrors.BaseErr{
				File:   "database/category_impl.go",
				Method: "UpdateCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
