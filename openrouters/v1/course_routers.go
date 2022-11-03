package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"opencourse/common"
	"strconv"
)

func (ctx *RouteContext) GetCourses(writer http.ResponseWriter, request *http.Request) {

	categoryId := chi.URLParam(request, "categoryId")

	urlValues := request.URL.Query()

	take := 5
	skip := 0
	var err error = nil

	if urlValues.Has("take") {
		take, err = strconv.Atoi(urlValues.Get("take"))
		if err != nil {
			writer.WriteHeader(400)
		}
	}

	if urlValues.Has("skip") {
		skip, err = strconv.Atoi(urlValues.Get("skip"))
		if err != nil {
			writer.WriteHeader(400)
		}
	}

	courses, err := ctx.DbContext.GetCourses(categoryId, int64(take), int64(skip))
	response := &OpenResponse[[]*common.Course]{Payload: courses}

	err = render.Render(writer, request, response)

	if err != nil {
		writer.WriteHeader(400)
	}
}
