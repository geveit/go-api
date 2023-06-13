package main

import (
	"database/sql"
	"log"

	"github.com/geveit/go-api/src/config"
	"github.com/geveit/go-api/src/item"
	"github.com/geveit/go-api/src/lib"
	"github.com/geveit/go-api/src/server"
	_ "github.com/jackc/pgx/v5/stdlib"
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

func getDatabaseConnection() (*sql.DB, error) {
	config, err := config.GetDatabaseConfig()
	if err != nil {
		return nil, err
	}

	return sql.Open("pgx", "postgres://"+config.Username+":"+config.Password+"@"+config.Host+":5432/"+config.Database)
}
