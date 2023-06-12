package item

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/geveit/go-api/src/server"
)

type Handler struct {
	repository Repository
}

func NewHandler(r Repository) *Handler {
	return &Handler{repository: r}
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	dbItem, err := h.getItem(w, r)
	if err != nil {
		return
	}

	server.WriteJSON(w, http.StatusOK, h.newItemResponse(dbItem))
}

func (h *Handler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.repository.GetAll()
	if err != nil {
		log.Printf("Error: %s", err.Error())
		server.WriteError(w, http.StatusInternalServerError, "Error fetching items")
		return
	}

	var response []*ItemResponse
	for _, item := range items {
		response = append(response, h.newItemResponse(item))
	}

	server.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	dbItem, err := h.getItem(w, r)
	if err != nil {
		return
	}

	if err := h.repository.Delete(dbItem.ID); err != nil {
		log.Printf("Error: %s", err.Error())
		server.WriteError(w, http.StatusInternalServerError, "Error deleting item")
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	itemRequest, err := h.decodeItemRequest(w, r)
	defer r.Body.Close()
	if err != nil {
		return
	}

	item := h.newItem(itemRequest)
	newId, err := h.repository.Insert(item)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		server.WriteError(w, http.StatusInternalServerError, "Error creating item")
		return
	}

	item.ID = newId

	server.WriteJSON(w, http.StatusCreated, h.newItemResponse(item))
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	itemRequest, err := h.decodeItemRequest(w, r)
	defer r.Body.Close()
	if err != nil {
		return
	}

	dbItem, err := h.getItem(w, r)
	if err != nil {
		return
	}

	item := h.newItem(itemRequest)
	item.ID = dbItem.ID

	if err := h.repository.Update(item); err != nil {
		log.Printf("Error: %s", err.Error())
		server.WriteError(w, http.StatusInternalServerError, "Error updating item")
		return
	}

	server.WriteJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) decodeItemRequest(w http.ResponseWriter, r *http.Request) (*ItemRequest, error) {
	itemRequest := ItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(&itemRequest); err != nil {
		server.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return nil, err
	}

	return &itemRequest, nil
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request) (*Item, error) {
	itemId, err := server.GetIdFromParam(r, "itemId")
	if err != nil {
		server.WriteError(w, http.StatusBadRequest, "Invalid id format")
		return nil, err
	}

	item, err := h.repository.Get(itemId)
	if err == ErrItemNotFound {
		server.WriteError(w, http.StatusNotFound, "Not found")
		return nil, err
	} else if err != nil {
		log.Printf("Error: %s", err.Error())
		server.WriteError(w, http.StatusInternalServerError, "Error fetching item")
		return nil, err
	}

	return item, nil
}

func (h *Handler) newItem(request *ItemRequest) *Item {
	return &Item{Name: request.Name}
}

func (h *Handler) newItemResponse(item *Item) *ItemResponse {
	return &ItemResponse{ID: item.ID, Name: item.Name}
}
