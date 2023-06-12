package lib

import (
	"database/sql"
)

type DBExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) (*sql.Row, error)
}

type sqlDBExecutor struct {
	db *sql.DB
}

func NewDBExecutor(db *sql.DB) DBExecutor {
	return &sqlDBExecutor{
		db: db,
	}
}

func (d *sqlDBExecutor) Exec(query string, args ...any) (sql.Result, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *sqlDBExecutor) Query(query string, args ...any) (*sql.Rows, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (d *sqlDBExecutor) QueryRow(query string, args ...any) (*sql.Row, error) {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}
