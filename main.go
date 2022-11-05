package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"opencourse/database"
	v1 "opencourse/openrouters/v1"
	"os"
)

func main() {

	// Loading environments
	conStr := os.Getenv("OPENCOURSE_CON_STR")
	sign := os.Getenv("OPENCOURSE_SIGN")
	smtpAccount := os.Getenv("OPENCOURSE_SMTP_ACCOUNT")
	smtpAccountPass := os.Getenv("OPENCOURSE_SMTP_ACCOUNT_PASS")
	baseEndpoint := os.Getenv("OPENCOURSE_ENDPOINT")

	dbContext := database.DbContext{}
	dbContext.Defaults(conStr, smtpAccount, smtpAccountPass, baseEndpoint)

	tokenAuth := jwtauth.New("HS256", []byte(sign), nil)

	err := dbContext.Connect()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := dbContext.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	logger := httplog.NewLogger("openlog", httplog.Options{
		JSON:    true,
		Concise: true,
	})

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Mount("/v1", v1.RouteTable(dbContext, tokenAuth))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Welcome to OpenCourses REST API"))
	})

	_ = http.ListenAndServe(":3000", r)
}
