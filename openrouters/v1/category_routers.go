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
		WriteErrResponse(writer, request, err, "can't get categories", 400)
		return
	}

	WriteResponse[[]*common.Category](writer, request, &categories)
}

func (ctx *RouteContext) PostCategory(writer http.ResponseWriter, request *http.Request) {

	openRequest := &OpenRequest[common.AddCategoryQuery]{}

	err := render.Bind(request, openRequest)

	if err != nil {
		WriteErrResponse(writer, request, err, "invalid model", 400)
		return
	}

	categoryId, err := ctx.DbContext.AddCategory(&openRequest.Payload)

	if err != nil {
		WriteErrResponse(writer, request, err, "can't create category", 400)
		return
	}

	WriteResponse[string](writer, request, &categoryId)
}
