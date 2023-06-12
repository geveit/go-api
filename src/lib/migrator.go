package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Migrator struct {
	dbExecutor DBExecutor
}

func NewMigrator(dbExecutor DBExecutor) *Migrator {
	return &Migrator{dbExecutor: dbExecutor}
}

func (m *Migrator) Migrate() error {
	createMigrationTableQuery := "CREATE TABLE IF NOT EXISTS migrations (file_name varchar primary key, \"timestamp\" timestamp)"
	if _, err := m.dbExecutor.Exec(createMigrationTableQuery); err != nil {
		return fmt.Errorf("Error creating migrations table: %w", err)
	}

	folderPath := "./migrations"
	folder, err := os.Open(folderPath)
	if err != nil {
		return fmt.Errorf("Couldn't open migrations folder: %w", err)
	}
	defer folder.Close()

	fileNames, err := folder.Readdirnames(0)
	if err != nil {
		return fmt.Errorf("Error reading files from migrations folder: %w", err)
	}

	rows, err := m.dbExecutor.Query("SELECT file_name FROM migrations")
	if err != nil {
		return fmt.Errorf("Error executing SELECT query on migrations table: %w", err)
	}
	defer rows.Close()

	filesMigrated := make(map[string]bool)
	for rows.Next() {
		var fileName string
		if err := rows.Scan(&fileName); err != nil {
			return fmt.Errorf("Error scanning filename from rows: %w", err)
		}

		filesMigrated[fileName] = true
	}

	for _, fileName := range fileNames {
		filePath := path.Join(folderPath, fileName)

		if filesMigrated[fileName] {
			continue
		}

		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("Error reading file %s: %w", fileName, err)
		}

		query := string(data)
		if _, err := m.dbExecutor.Exec(query); err != nil {
			return fmt.Errorf("Error migrating %s: %w", fileName, err)
		}

		insertQuery := "INSERT INTO migrations (file_name, \"timestamp\") VALUES ($1, now())"
		if _, err := m.dbExecutor.Exec(insertQuery, fileName); err != nil {
			return fmt.Errorf("Error saving migration %s: %w", fileName, err)
		}

		log.Printf("Successfully migrated %s", fileName)
	}

	return nil
}
