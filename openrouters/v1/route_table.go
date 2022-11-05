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

		r.Get("/courses/{categoryId}", rtx.GetCourses)
		r.Get("/categories/{lang}", rtx.GetCategories)
		r.Post("/categories", rtx.PostCategory)
	})

	r.Group(func(r chi.Router) {
		r.Post("/auth/signin", rtx.Login)
		r.Post("/auth/signup", nil)
	})

	return r
}
