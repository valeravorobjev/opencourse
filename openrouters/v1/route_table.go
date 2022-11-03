package v1

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"opencourse/database"
)

func RouteTable(dbContext database.DbContext) http.Handler {
	r := chi.NewRouter()
	rtx := RouteContext{DbContext: dbContext}

	// Courses group
	r.Group(func(r chi.Router) {
		r.Get("/courses/{categoryId}", rtx.GetCourses)
	})

	r.Group(func(r chi.Router) {
		r.Get("/categories/{lang}", rtx.GetCategories)
		r.Post("/categories", rtx.PostCategory)
	})

	return r
}
