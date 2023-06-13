package main

import (
	"log"

	"github.com/geveit/go-api/src/item"
	"github.com/geveit/go-api/src/lib"
	"github.com/geveit/go-api/src/server"
)

func main() {
	server := server.NewServer(":3000")

	db, err := getDatabaseConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}
	defer db.Close()

	dbExecutor := lib.NewDBExecutor(db)

	migrator := lib.NewMigrator(dbExecutor)
	if err := migrator.Migrate(); err != nil {
		log.Fatal(err)
	}

	itemRepository := item.NewRepository(dbExecutor)
	itemHandler := item.NewHandler(itemRepository)

	item.RegisterRoutes(server.Router, itemHandler)

	server.Run()
}
