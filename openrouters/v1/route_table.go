package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"opencourse/database"
)

func RouteTable(dbContext database.DbContext, tokenAuth *jwtauth.JWTAuth) http.Handler {
	r := chi.NewRouter()
	rtx := RouteContext{DbContext: dbContext, TokenAuth: tokenAuth}

	r.Group(func(r chi.Router) {

		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/courses/{categoryId}/list", rtx.GetCourses)
		r.Get("/courses/{courseId}", rtx.GetCourses)
		r.Post("/courses", rtx.PostCourse)

		r.Get("/stages/{courseId}/list", rtx.GetStages)
		r.Get("/stages/{stageId}", rtx.GetStage)
		r.Post("/stages", rtx.PostStage)
		r.Put("/stages", rtx.PutStage)

		r.Get("/categories/{lang}", rtx.GetCategories)
		r.Post("/categories", rtx.PostCategory)
	})

	r.Group(func(r chi.Router) {
		r.Post("/auth/login", rtx.Login)
		r.Post("/auth/register", rtx.Register)
		r.Get("/auth/confirm/{id}/{code}", rtx.Confirm)
	})

	return r
}
