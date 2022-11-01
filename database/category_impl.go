package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
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
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
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
			return nil, openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "database/mongodb/category_impl.go",
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

	for _, openSubCategory := range addCategoryQuery.SubCategories {

		subCategory := &DbSubCategory{Number: openSubCategory.Number, Name: openSubCategory.Name}
		category.SubCategories = append(category.SubCategories, subCategory)
	}

	result, err := col.InsertOne(context.Background(), category)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
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
UpdateCategory update category name
@cid - category id
@names - category names
*/
func (ctx *DbContext) UpdateCategory(cid string, name string) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	filter := bson.D{
		{"_id", cid},
	}

	update := bson.D{
		{
			"$set", bson.D{{"name", name}},
		},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
				Method: "UpdateCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
UpdateSubCategory update sub category names
@cid - category id
@scn - sub category number
@names - sub category names
*/
func (ctx *DbContext) UpdateSubCategory(cid string, scn int, name string) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"_id", cid}},
			bson.D{{"sub_categories.number", scn}},
		}},
	}

	update := bson.D{
		{
			"$set", bson.D{{"sub_categories.$.name", name}},
		},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
				Method: "UpdateSubCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}

/*
AddSubCategory add sub category for category.
@cid   - category id
@name - sub category name
@lang - sub category language
*/
func (ctx *DbContext) AddSubCategory(cid string, name string) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	categoryId, err := primitive.ObjectIDFromHex(cid)

	if err != nil {
		return openerrors.OpenInvalidIdErr{
			Default: openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "database/mongodb/category_impl.go",
					Method: "AddSubCategory",
				},
				Msg: fmt.Sprintf("can't convert cid %s to ObjectID with method primitive.ObjectIDFromHex(cid)", cid),
			},
			Id:        cid,
			Converter: " primitive.ObjectIDFromHex(cid)",
		}
	}

	filter := bson.D{
		{"_id", categoryId},
	}

	var category DbCategory
	err = col.FindOne(context.Background(), filter).Decode(&category)

	if err != nil && err != mongo.ErrNoDocuments {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
				Method: "AddSubCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var subCategory *DbSubCategory

	if len(category.SubCategories) == 0 {
		subCategory = &DbSubCategory{Number: 0, Name: name}
	} else {
		lastSubCategory := category.SubCategories[len(category.SubCategories)-1]

		subCategory =
			&DbSubCategory{Number: lastSubCategory.Number + 1, Name: name}
	}

	update := bson.D{
		{
			"$push", bson.D{{"sub_categories", subCategory}},
		},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
				Method: "AddSubCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil

}

/*
DeleteSubCategory delete sub category from category
@cid - category id
@scn - sub category number
*/
func (ctx *DbContext) DeleteSubCategory(cid string, scn int) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	categoryId, err := primitive.ObjectIDFromHex(cid)

	if err != nil {
		return openerrors.OpenInvalidIdErr{
			Default: openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "database/mongodb/category_impl.go",
					Method: "DeleteSubCategory",
				},
				Msg: fmt.Sprintf("can't convert cid %s to ObjectID with method primitive.ObjectIDFromHex(cid)", cid),
			},
			Id:        cid,
			Converter: "primitive.ObjectIDFromHex(cid)",
		}
	}

	filter := bson.D{
		{"_id", categoryId},
	}

	var category DbCategory
	err = col.FindOne(context.Background(), filter).Decode(&category)

	if err != nil && err != mongo.ErrNoDocuments {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
				Method: "DeleteSubCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	if len(category.SubCategories) == 0 {
		return nil
	}

	var subCategories []*DbSubCategory
	number := 0
	for _, item := range category.SubCategories {
		if item.Number == scn {
			continue
		}
		item.Number = number
		subCategories = append(subCategories, item)
		number++
	}

	update := bson.D{
		{
			"$set", bson.D{{"sub_categories", subCategories}},
		},
	}

	_, err = col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "database/mongodb/category_impl.go",
				Method: "DeleteSubCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
