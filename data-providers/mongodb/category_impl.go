package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"opencourse/common/openerrors"
)

// GetCategories return all categories from db
func (ctx *ContextMongoDb) GetCategories() ([]*Category, error) {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	cursor, err := col.Find(context.Background(), bson.D{})

	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/category_impl.go",
				Method: "GetCategories",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var categories []*Category

	err = cursor.Decode(categories)
	if err != nil {
		return nil, openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/category_impl.go",
				Method: "GetCategories",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return categories, nil
}

// AddCategory method for add new category and sub categories to database
func (ctx *ContextMongoDb) AddCategory(category *Category) (string, error) {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	result, err := col.InsertOne(context.Background(), category)

	if err != nil {
		return "", openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/category_impl.go",
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
func (ctx *ContextMongoDb) UpdateCategory(cid string, names []*GlobStr) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	filter := bson.D{
		{"_id", cid},
	}

	update := bson.D{
		{
			"$set", bson.D{{"names", names}},
		},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/category_impl.go",
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
func (ctx *ContextMongoDb) UpdateSubCategory(cid string, scn int, names []*GlobStr) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	filter := bson.D{
		{"$and", bson.A{
			bson.D{{"_id", cid}},
			bson.D{{"sub_categories.number", scn}},
		}},
	}

	update := bson.D{
		{
			"$set", bson.D{{"sub_categories.$.names", names}},
		},
	}

	_, err := col.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/category_impl.go",
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
func (ctx *ContextMongoDb) AddSubCategory(cid string, name string, lang string) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	categoryId, err := primitive.ObjectIDFromHex(cid)

	if err != nil {
		return openerrors.OpenInvalidIdErr{
			Default: openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/category_impl.go",
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

	var category Category
	err = col.FindOne(context.Background(), filter).Decode(&category)

	if err != nil && err != mongo.ErrNoDocuments {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/category_impl.go",
				Method: "AddSubCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	var subCategory *SubCategory

	if len(category.SubCategories) == 0 {
		subCategory = &SubCategory{Number: 0, Names: []*GlobStr{
			{
				Lang: lang,
				Text: name,
			},
		}}
	} else {
		lastSubCategory := category.SubCategories[len(category.SubCategories)-1]

		subCategory =
			&SubCategory{Number: lastSubCategory.Number + 1, Names: []*GlobStr{
				{
					Lang: lang,
					Text: name,
				},
			}}
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
				File:   "data-providers/mongodb/category_impl.go",
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
func (ctx *ContextMongoDb) DeleteSubCategory(cid string, scn int) error {
	col := ctx.Client.Database(DbName).Collection(CategoryCollection)

	categoryId, err := primitive.ObjectIDFromHex(cid)

	if err != nil {
		return openerrors.OpenInvalidIdErr{
			Default: openerrors.OpenDefaultErr{
				BaseErr: openerrors.OpenBaseErr{
					File:   "data-providers/mongodb/category_impl.go",
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

	var category Category
	err = col.FindOne(context.Background(), filter).Decode(&category)

	if err != nil && err != mongo.ErrNoDocuments {
		return openerrors.OpenDbErr{
			BaseErr: openerrors.OpenBaseErr{
				File:   "data-providers/mongodb/category_impl.go",
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

	var subCategories []*SubCategory
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
				File:   "data-providers/mongodb/category_impl.go",
				Method: "DeleteSubCategory",
			},
			DbName: ctx.DbName,
			ConStr: ctx.Uri,
			DbErr:  err.Error(),
		}
	}

	return nil
}
