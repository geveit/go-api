package config

import (
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
	data, err := getJsonReader().readConfigJson()
	if err != nil {
		return nil, err
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
