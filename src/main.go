package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/geveit/go-api/src/item"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const port = ":3000"

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	db, err := sql.Open("pgx", "postgres://postgres:super_secret@localhost:5432/go-api")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	itemRepository := item.NewRepository(db)
	itemHandler := item.NewHandler(itemRepository)

	item.RegisterRoutes(r, itemHandler)

	log.Fatal(http.ListenAndServe(":3000", r))
}
