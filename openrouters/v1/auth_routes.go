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
		WriteErrResponse[string](writer, request, "invalid model", 400)
		return
	}

	_, tokenString, err := ctx.TokenAuth.Encode(
		map[string]interface{}{
			"login": openRequest.Payload.Login,
			"exp":   time.Now().Add(time.Minute * 5).Unix(),
		})

	if err != nil {
		WriteErrResponse[string](writer, request, "create token error", 400)
		return
	}

	WriteResponse[string](writer, request, &tokenString)
}
