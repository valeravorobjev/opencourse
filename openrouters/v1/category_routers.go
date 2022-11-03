package v1

import (
	"github.com/go-chi/render"
	"net/http"
	"opencourse/common"
)

func (ctx *RouteContext) GetCategories(writer http.ResponseWriter, request *http.Request) {

	categories, err := ctx.DbContext.GetCategories("en")
	response := &OpenResponse[[]*common.Category]{Data: categories}

	err = render.Render(writer, request, response)

	if err != nil {
		writer.WriteHeader(400)
	}
}
