package querier

import (
	"database/sql"
)

func Exec(db *sql.DB, query string, args ...any) (sql.Result, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Query(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func QueryRow(db *sql.DB, query string, args ...any) (*sql.Row, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args), nil
}
