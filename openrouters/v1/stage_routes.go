package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"opencourse/common"
	"strconv"
)

func (ctx *RouteContext) GetStages(writer http.ResponseWriter, request *http.Request) {

	courseId := chi.URLParam(request, "courseId")

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

	stagePreviews, err := ctx.DbContext.GetStages(courseId, int64(take), int64(skip))

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't get stages."}, 400)
		return
	}

	WriteResponse[[]*common.StagePreview](writer, request, &stagePreviews)
}

func (ctx *RouteContext) GetStage(writer http.ResponseWriter, request *http.Request) {

	stageId := chi.URLParam(request, "stageId")

	stage, err := ctx.DbContext.GetStage(stageId)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't get stage."}, 400)
		return
	}

	WriteResponse[common.Stage](writer, request, stage)

}

func (ctx *RouteContext) PostStage(writer http.ResponseWriter, request *http.Request) {
	// Check user role. If user is not in role, return.
	ok := InRole(writer, request, common.RoleAdmin)
	if !ok {
		return
	}

	openRequest := &Request[common.AddStageQuery]{}

	err := render.Bind(request, openRequest)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrBinding, Message: "Invalid model"}, 400)
		return
	}

	id, err := ctx.DbContext.AddStage(&openRequest.Payload)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't add stage."}, 400)
		return
	}

	WriteResponse[string](writer, request, &id)
}

func (ctx *RouteContext) PutStage(writer http.ResponseWriter, request *http.Request) {
	// Check user role. If user is not in role, return.
	ok := InRole(writer, request, common.RoleAdmin)
	if !ok {
		return
	}

	openRequest := &Request[common.UpdateStageQuery]{}

	err := render.Bind(request, openRequest)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrBinding, Message: "Invalid model"}, 400)
		return
	}

	err = ctx.DbContext.UpdateStage(&openRequest.Payload)

	if err != nil {
		WriteErrResponse(writer, request, err,
			&ResponseError{Code: ErrInternal, Message: "Internal error. Can't update stage."}, 400)
		return
	}

	result := "success"
	WriteResponse[string](writer, request, &result)
}
