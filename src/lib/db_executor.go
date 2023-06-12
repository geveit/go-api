package lib

import (
	"database/sql"
	"fmt"
)

type DBExecutor interface {
	Begin() error
	Commit() error
	Rollback() error
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) (*sql.Row, error)
}

type sqlDBExecutor struct {
	db *sql.DB
	tx *sql.Tx
}

func NewDBExecutor(db *sql.DB) DBExecutor {
	return &sqlDBExecutor{
		db: db,
	}
}

func (d *sqlDBExecutor) Begin() error {
	tx, err := d.db.Begin()
	d.tx = tx
	return err
}

func (d *sqlDBExecutor) Commit() error {
	if d.tx != nil {
		err := d.tx.Commit()
		d.tx = nil
		return err
	}

	return fmt.Errorf("Cannot commit, transaction not started")
}

func (d *sqlDBExecutor) Rollback() error {
	if d.tx != nil {
		err := d.tx.Rollback()
		d.tx = nil
		return err
	}

	return nil
}

func (d *sqlDBExecutor) Exec(query string, args ...any) (sql.Result, error) {
	stmt, err := d.prepare(query)
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
	stmt, err := d.prepare(query)
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
	stmt, err := d.prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRow(args...), nil
}

func (d *sqlDBExecutor) prepare(query string) (*sql.Stmt, error) {
	if d.tx != nil {
		return d.tx.Prepare(query)
	}

	return d.db.Prepare(query)
}
