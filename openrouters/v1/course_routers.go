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
		WriteErrResponse(writer, request, err, &ResponseError{Code: ErrParameter, Message: "Wrong take parameter."}, 400)
		return
	}

	if urlValues.Has("skip") {
		skip, err = strconv.Atoi(urlValues.Get("skip"))

		WriteErrResponse(writer, request, err, &ResponseError{Code: ErrParameter, Message: "Wrong skip parameter."}, 400)
		return
	}

	courses, err := ctx.DbContext.GetCourses(categoryId, int64(take), int64(skip))

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't get courses."}, 400)
		return
	}

	WriteResponse[[]*common.Course](writer, request, &courses)
}

func (ctx *RouteContext) GetCourse(writer http.ResponseWriter, request *http.Request) {

	courseId := chi.URLParam(request, "courseId")

	course, err := ctx.DbContext.GetCourse(courseId)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't get course."}, 400)
		return
	}

	WriteResponse[common.Course](writer, request, course)

}

func (ctx *RouteContext) PostCourse(writer http.ResponseWriter, request *http.Request) {
	// Check user role. If user is not in role, return.
	ok := InRole(writer, request, common.RoleAdmin)
	if !ok {
		return
	}

	openRequest := &Request[common.AddCourseQuery]{}

	err := render.Bind(request, openRequest)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrBinding, Message: "Invalid model"}, 400)
		return
	}

	id, err := ctx.DbContext.AddCourse(&openRequest.Payload)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't add course."}, 400)
		return
	}

	WriteResponse[string](writer, request, &id)
}
