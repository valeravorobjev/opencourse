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
	response := &OpenResponse[[]*common.Category]{Payload: categories}

	err = render.Render(writer, request, response)

	if err != nil {
		writer.WriteHeader(400)
	}
}

func (ctx *RouteContext) PostCategory(writer http.ResponseWriter, request *http.Request) {

	openRequest := &OpenRequest[common.AddCategoryQuery]{}

	err := render.Bind(request, openRequest)

	categoryId, err := ctx.DbContext.AddCategory(&openRequest.Payload)

	response := &OpenResponse[string]{Payload: categoryId}

	err = render.Render(writer, request, response)

	if err != nil {
		writer.WriteHeader(400)
	}
}
