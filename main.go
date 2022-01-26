package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/newtonmunene99/go-books/ent"
	"github.com/newtonmunene99/go-books/ent/migrate"
	"github.com/newtonmunene99/go-books/env"
)

var client *ent.Client

func Open(databaseUrl string) *ent.Client {
	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Create an ent.Driver from `db`.
	drv := entsql.OpenDB(dialect.Postgres, db)
	return ent.NewClient(ent.Driver(drv))
}

func main() {

	client = Open(fmt.Sprintf("postgresql://%v:%v@127.0.0.1/go-books", env.DB_USERNAME, env.DB_PASSWORD))

	ctx := context.Background()
	err := client.Schema.Create(
		ctx,
		migrate.WithDropIndex(true),
		migrate.WithDropColumn(true),
	)

	if err != nil {
		log.Fatal(err)
	}

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()

	r.Path("/").Methods("GET").HandlerFunc(ListBooksHandler)
	r.Path("/").Methods("POST").HandlerFunc(CreateBookHandler)

	r.Path("/categories").Methods("GET").HandlerFunc(ListCategoriesHandler)
	r.Path("/categories").Methods("POST").HandlerFunc(CreateCategoryHandler)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {

		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}

	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt)

	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	srv.Shutdown(ctx)
	log.Println("shutting down")

	os.Exit(0)

}

func ListCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	categories, err := client.Category.Query().All(context.Background())

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("%v", categories)))

}

func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {

	var body ent.Category

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := client.Category.
		Create().
		SetName(body.Name).
		Save(context.Background())

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("%v", category)))

}

func ListBooksHandler(w http.ResponseWriter, r *http.Request) {

	books, err := client.Book.Query().All(context.Background())

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("%v", books)))

}

func CreateBookHandler(w http.ResponseWriter, r *http.Request) {

	var body ent.Book

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	category, err := client.Category.Get(context.Background(), body.CategoryID)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("Category with Id %v Not Found: %v", body.CategoryID, err)))
	}

	book, err := client.Book.
		Create().
		SetTitle(body.Title).
		SetYear(body.Year).
		SetAuthor(body.Author).
		SetCategoryID(category.ID).
		AddCategory(category).
		Save(context.Background())

	if err != nil {

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("%v", book)))

}
