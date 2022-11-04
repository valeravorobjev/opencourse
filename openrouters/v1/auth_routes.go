package v1

import (
	"github.com/go-chi/render"
	"net/http"
	"opencourse/common"
	"time"
)

func (ctx *RouteContext) Login(writer http.ResponseWriter, request *http.Request) {

	openRequest := &OpenRequest[common.LoginQuery]{}

	err := render.Bind(request, openRequest)
	if err != nil {
		WriteErrResponse(writer, request, err, "invalid model", 400)
		return
	}

	_, tokenString, err := ctx.TokenAuth.Encode(
		map[string]interface{}{
			"login": openRequest.Payload.Login,
			"exp":   time.Now().Add(time.Minute * 60).Unix(),
		})

	if err != nil {
		WriteErrResponse(writer, request, err, "create token error", 400)
		return
	}

	WriteResponse[string](writer, request, &tokenString)
}
