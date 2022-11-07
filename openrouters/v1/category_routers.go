package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"opencourse/common"
)

func (ctx *RouteContext) GetCategories(writer http.ResponseWriter, request *http.Request) {

	lang := chi.URLParam(request, "lang")

	categories, err := ctx.DbContext.GetCategories(lang)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't get categories."}, 400)
		return
	}

	WriteResponse[[]*common.Category](writer, request, &categories)
}

func (ctx *RouteContext) PostCategory(writer http.ResponseWriter, request *http.Request) {

	openRequest := &Request[common.AddCategoryQuery]{}

	err := render.Bind(request, openRequest)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrBinding, Message: "Invalid model"}, 400)
		return
	}

	categoryId, err := ctx.DbContext.AddCategory(&openRequest.Payload)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrUserAlreadyExists, Message: "Internal error. Can't create category."}, 400)
		return
	}

	WriteResponse[string](writer, request, &categoryId)
}

func (ctx *RouteContext) PutCategory(writer http.ResponseWriter, request *http.Request) {
	openRequest := &Request[common.UpdateCategoryQuery]{}

	err := render.Bind(request, openRequest)

	if err != nil {
		WriteErrResponse(writer, request, err, &ResponseError{Code: ErrBinding, Message: "Invalid model"}, 400)
		return
	}

	err = ctx.DbContext.UpdateCategory(openRequest.Payload.CategoryId, openRequest.Payload.Name, openRequest.Payload.Lang)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't update category."}, 400)
		return
	}

	result := "success"
	WriteResponse[string](writer, request, &result)
}
