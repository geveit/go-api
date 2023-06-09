package item

import (
	"database/sql"
	"errors"

	"github.com/geveit/go-api/src/lib/querier"
)

var (
	ErrItemNotFound = errors.New("item not found")
)

type Repository interface {
	Insert(item *Item) (uint, error)
	Delete(id uint) error
	Update(item *Item) error
	Get(id uint) (*Item, error)
	GetAll() ([]*Item, error)
}

type sqlRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqlRepository{db: db}
}

func (r *sqlRepository) Insert(item *Item) (uint, error) {
	query := "INSERT INTO items (name) values (?)"

	result, err := querier.Exec(r.db, query, item.Name)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (r *sqlRepository) Delete(id uint) error {
	query := "DELETE FROM items WHERE id = ?"

	if _, err := querier.Exec(r.db, query, id); err != nil {
		return err
	}

	return nil
}

func (r *sqlRepository) Update(item *Item) error {
	query := "UPDATE items SET name = ? WHERE id = ?"

	if _, err := querier.Exec(r.db, query, item.Name, item.ID); err != nil {
		return err
	}

	return nil
}

func (r *sqlRepository) Get(id uint) (*Item, error) {
	query := "SELECT id, name FROM items WHERE id = ?"

	row, err := querier.QueryRow(r.db, query, id)
	if err != nil {
		return nil, err
	}

	var item Item
	err = row.Scan(&item.ID, &item.Name)

	if err == sql.ErrNoRows {
		return nil, ErrItemNotFound
	} else if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *sqlRepository) GetAll() ([]*Item, error) {
	query := "SELECT id, name FROM users"
	rows, err := querier.Query(r.db, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.ID, &item.Name)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
