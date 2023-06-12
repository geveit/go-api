package main

import (
	"database/sql"
	"log"

	"github.com/geveit/go-api/src/item"
	"github.com/geveit/go-api/src/lib"
	"github.com/geveit/go-api/src/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	server := server.NewServer(":3000")

	db, err := sql.Open("pgx", "postgres://postgres:super_secret@localhost:5432/go-api")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	querier := lib.NewQuerier(db)

	itemRepository := item.NewRepository(querier)
	itemHandler := item.NewHandler(itemRepository)

	item.RegisterRoutes(server.Router, itemHandler)

	server.Run()
}
