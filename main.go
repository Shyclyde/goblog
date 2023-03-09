package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

var (
	router *chi.Mux
	db     *sql.DB
)

type Article struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	router = chi.NewRouter()
	router.Use(middleware.Recoverer)

	err := dbConnect()
	catch(err)
	defer db.Close()

	router.Use(ChangeMethod)
	router.Get("/", GetAllArticles)
	router.Route("/articles",
		func(r chi.Router) {
			r.Get("/", NewArticle)
			r.Post("/", CreateArticle)
			r.Route("/{articleID}",
				func(r chi.Router) {
					r.Use(ArticleCtx)
					r.Get("/", GetArticle)
					r.Put("/", UpdateArticle)
					r.Delete("/", DeleteArticle)
					r.Get("/edit", EditArticle)
				})
		})

	err = http.ListenAndServe(":8005", router)
	catch(err)

	// article := &Article{0, "test", "testing"}
	// err = dbCreateArticle(db, article)
	// catch(err)

	// article, err := dbGetArticle(db, "3")
	// catch(err)
	// fmt.Println(article.Content)

	// article := &Article{4, "yay2", "whoasdfasdfopie"}
	// updated, err := dbUpdateArticle(db, article)
	// catch(err)
	// if updated {
	// 	fmt.Println("Updated!")
	// } else {
	// 	fmt.Println("No rows updated")
	// }

	// deleted, err := dbDeleteArticle(db, "3")
	// catch(err)
	// if deleted {
	// 	fmt.Println("Article deleted!")
	// } else {
	// 	fmt.Println("No articles deleted.")
	// }

	articles, err := dbGetAllArticles()
	catch(err)
	for _, item := range articles {
		fmt.Printf("ID: %d, Title: %s, Content: %s\n", item.ID, item.Title, item.Content)
	}
}

func ChangeMethod(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				switch method := r.PostFormValue("_method"); method {
				case http.MethodPut:
					fallthrough
				case http.MethodPatch:
					fallthrough
				case http.MethodDelete:
					r.Method = method
				default:
				}
			}
			next.ServeHTTP(w, r)
		},
	)
}

func ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			articleID := chi.URLParam(r, "articleID")
			article, err := dbGetArticle(articleID)
			if err != nil {
				fmt.Println(err)
				http.Error(w, http.StatusText(404), 404)
				return
			}
			ctx := context.WithValue(r.Context(), "article", article)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func GetAllArticles(next http.Handler) http.Handler {
	return nil
}
