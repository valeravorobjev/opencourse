package v1

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"gopkg.in/gomail.v2"
	"net/http"
	"opencourse/common"
	"time"
)

// Login TODO: this route is not complete!
func (ctx *RouteContext) Login(writer http.ResponseWriter, request *http.Request) {

	openRequest := &OpenRequest[common.LoginQuery]{}

	err := render.Bind(request, openRequest)
	if err != nil {
		WriteErrResponse(writer, request, err, "invalid model", 400)
		return
	}

	// TODO: put here user logic

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

// Register route
func (ctx *RouteContext) Register(writer http.ResponseWriter, request *http.Request) {
	openRequest := &OpenRequest[common.RegisterQuery]{}

	err := render.Bind(request, openRequest)
	if err != nil {
		WriteErrResponse(writer, request, err, "invalid model", 400)
		return
	}

	user, err := ctx.DbContext.GetUserByLogin(openRequest.Payload.Login)

	if err != nil {
		WriteErrResponse(writer, request, err, "registration error", 400)
		return
	}

	if user != nil {
		WriteErrResponse(writer, request, nil,
			fmt.Sprintf("user with login %s already exists", openRequest.Payload.Login), 400)
		return
	}

	userConfirm, err := ctx.DbContext.GetUserConfirmByLogin(openRequest.Payload.Login)

	if err != nil {
		WriteErrResponse(writer, request, err, "registration error", 400)
		return
	}

	// if user don't exist in collection users but user_confirms user has confirmed
	if userConfirm != nil && userConfirm.Confirmed == true {
		WriteErrResponse(writer, request, err, "user is incorrect, registration error, please contact to support", 400)
		return
	}

	if userConfirm != nil {
		WriteErrResponse(writer, request, err, "user with this login already exists", 400)
		return
	}

	userConfirm, err = ctx.DbContext.AddUserConfirm(&openRequest.Payload)

	if err != nil {
		WriteErrResponse(writer, request, err, "registration error", 400)
		return
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", "confirm@opencourse.com")
	msg.SetHeader("To", userConfirm.Email)
	msg.SetHeader("Subject", "Confirm registration")

	link := fmt.Sprintf("%s/%s/%s/%s", ctx.DbContext.Endpoint, "v1/auth/confirm", userConfirm.Id, userConfirm.ConfirmaCode)
	text := `
<h3>OpenCourse confirmation of registration.</h3>
<p>This email is automatically sent by OpenCourse. Don't answer it.</p>
`
	msg.SetBody("text/html", fmt.Sprintf("%s <p>Please, follow the link to <a href='%s'>confirm</a></p>", text, link))

	n := gomail.NewDialer("smtp.gmail.com", 587, ctx.DbContext.SmtpAccount, ctx.DbContext.SmtpAccountPass)

	if err := n.DialAndSend(msg); err != nil {
		time.Sleep(time.Second * 1) // If error, try to send message with pause 1 sec

		if err := n.DialAndSend(msg); err != nil {

			tempErr := err

			err := ctx.DbContext.DeleteUserConfirm(userConfirm.Id)

			if err != nil {
				err = errors.New(fmt.Sprintf("%s | %s", tempErr, err))
			}

			WriteErrResponse(writer, request, err, "registration error", 400)
			return
		}
	}

	message := "Please confirm your registration. " +
		"The confirmation link has been sent to your email"
	WriteResponse[string](writer, request, &message)

}

func (ctx *RouteContext) Confirm(writer http.ResponseWriter, request *http.Request) {
	confirmId := chi.URLParam(request, "id")
	code := chi.URLParam(request, "code")

	userConfirm, err := ctx.DbContext.GetUserConfirm(confirmId)

	if err != nil {
		WriteErrResponse(writer, request, err, "Confirmation error. Try follow the link again or registration.", 400)
		return
	}

	if userConfirm.ConfirmaCode != code {
		WriteErrResponse(writer, request, errors.New("registration confirmation codes do not match"), "The link is not valid.", 400)
		return
	}

	if userConfirm.Confirmed == true {
		WriteErrResponse(writer, request, errors.New(fmt.Sprintf("user login = %s id = %s already confirmed", userConfirm.Login, userConfirm.Id)), "The link is not valid.", 400)
		return
	}

	user, err := ctx.DbContext.GetUserByLogin(userConfirm.Login)

	if err != nil {
		WriteErrResponse(writer, request, err, "Confirmation error. Try follow the link again or registration.", 400)
		return
	}

	if user != nil {
		WriteErrResponse(writer, request, errors.New("user already exist"), "The link is not valid.", 400)
		return
	}

	addUserQuery := common.AddUserQuery{
		Login:    userConfirm.Login,
		Password: userConfirm.Password,
		Email:    userConfirm.Email,
		Name:     userConfirm.Name,
		Avatar:   userConfirm.Avatar,
		Roles:    []string{common.RoleUser},
	}

	_, err = ctx.DbContext.AddUser(&addUserQuery)

	if err != nil {
		WriteErrResponse(writer, request, err, "Confirmation error. Try follow the link again or registration.", 400)
		return
	}

	err = ctx.DbContext.SetConfirmed(confirmId)

	if err != nil {
		WriteErrResponse(writer, request, err, "Confirmation error. Try follow the link again or registration.", 400)
		return
	}

	_, _ = writer.Write([]byte("<b>Registration SUCCESS completed!</b>"))
}
