package lib

import (
	"database/sql"
)

type Querier interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) (*sql.Row, error)
}

type sqlQuerier struct {
	db *sql.DB
}

func NewQuerier(db *sql.DB) Querier {
	return &sqlQuerier{
		db: db,
	}
}

func (q *sqlQuerier) Exec(query string, args ...any) (sql.Result, error) {
	stmt, err := q.db.Prepare(query)
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

func (q *sqlQuerier) Query(query string, args ...any) (*sql.Rows, error) {
	stmt, err := q.db.Prepare(query)
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

func (q *sqlQuerier) QueryRow(query string, args ...any) (*sql.Row, error) {
	stmt, err := q.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args), nil
}
