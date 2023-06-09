package item

import (
	"time"
)

type Item struct {
	ID   uint
	Name string
}

type ItemRequest struct {
	Name string `json:"name"`
}

type ItemResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
