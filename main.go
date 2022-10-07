package main

import (
	"net/http"
	"opencourse/data-providers/mongodb"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {

		db := &mongodb.ContextMongoDb{}

		// Init default values
		db.Defaults()

		err := db.Connect("mongodb://localhost")
		if err != nil {
			panic(err)
		}

		defer func() {
			err = db.Disconnect()
			if err != nil {
				panic(err)
			}
		}()

		//category := mongodb.Category{
		//	Names: []*mongodb.GlobStr{
		//		{
		//			Lang: mongodb.LangRu,
		//			Text: "Программирование",
		//		},
		//	},
		//	SubCategories: []*mongodb.SubCategory{
		//		{
		//			Number: 1,
		//			Names: []*mongodb.GlobStr{
		//				{
		//					Lang: mongodb.LangRu,
		//					Text: "Язык C# и платформа .NET",
		//				},
		//			},
		//		},
		//	},
		//}
		//
		//id, err := db.AddCategory(&category)
		//
		//if err != nil {
		//	panic(err)
		//}

		err = db.DeleteSubCategory("62fca089be2a1ee3caae2e82", 2)

		if err != nil {
			panic(err)
		}

		_, _ = w.Write([]byte("62fca089be2a1ee3caae2e82"))
	})
	_ = http.ListenAndServe(":3000", r)
}
