package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_USERNAME = "DATABASE_USERNAME"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_DATABASE = "DATABASE_DATABASE"
)

func getDatabaseConnection() (*sql.DB, error) {
	filePath := "../config.json"
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config.json: %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("Failed to unmarshall json data: %w", err)
	}

	config := data["database"].(map[string]any)
	host := config["host"].(string)
	database := config["database"].(string)
	password := config["password"].(string)
	username := config["username"].(string)

	if os.Getenv(DATABASE_HOST) != "" {
		host = os.Getenv(DATABASE_HOST)
	}
	if os.Getenv(DATABASE_DATABASE) != "" {
		database = os.Getenv(DATABASE_DATABASE)
	}
	if os.Getenv(DATABASE_USERNAME) != "" {
		username = os.Getenv(DATABASE_USERNAME)
	}
	if os.Getenv(DATABASE_PASSWORD) != "" {
		password = os.Getenv(DATABASE_PASSWORD)
	}

	return sql.Open("pgx", "postgres://"+username+":"+password+"@"+host+":5432/"+database)
}
