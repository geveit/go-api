package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Username string
	Password string
	Database string
}

const (
	DATABASE_HOST     = "DATABASE_HOST"
	DATABASE_USERNAME = "DATABASE_USERNAME"
	DATABASE_PASSWORD = "DATABASE_PASSWORD"
	DATABASE_DATABASE = "DATABASE_DATABASE"
)

func GetDatabaseConfig() (*DatabaseConfig, error) {
	filePath := "../../config.json"
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open config.json: %w", err)
	}

	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("Failed to unmarshall json data: %w", err)
	}

	config := DatabaseConfig{}

	configJson := data["database"].(map[string]any)
	config.Host = configJson["host"].(string)
	config.Database = configJson["database"].(string)
	config.Password = configJson["password"].(string)
	config.Username = configJson["username"].(string)

	if os.Getenv(DATABASE_HOST) != "" {
		config.Host = os.Getenv(DATABASE_HOST)
	}
	if os.Getenv(DATABASE_DATABASE) != "" {
		config.Database = os.Getenv(DATABASE_DATABASE)
	}
	if os.Getenv(DATABASE_USERNAME) != "" {
		config.Username = os.Getenv(DATABASE_USERNAME)
	}
	if os.Getenv(DATABASE_PASSWORD) != "" {
		config.Password = os.Getenv(DATABASE_PASSWORD)
	}

	return &config, nil
}
