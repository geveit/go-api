package item

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/geveit/go-api/src/helper"
	"github.com/go-chi/chi"
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

	helper.JsonResponse(w, http.StatusOK, h.newItemResponse(dbItem))
}

func (h *Handler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.repository.GetAll()
	if err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error fetching items")
		return
	}

	var response []*ItemResponse
	for _, item := range items {
		response = append(response, h.newItemResponse(item))
	}

	helper.JsonResponse(w, http.StatusOK, response)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	dbItem, err := h.getItem(w, r)
	if err != nil {
		return
	}

	if err := h.repository.Delete(dbItem.ID); err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error deleting item")
		return
	}

	helper.JsonResponse(w, http.StatusNoContent, nil)
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
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error creating item")
		return
	}

	item.ID = newId

	helper.JsonResponse(w, http.StatusCreated, h.newItemResponse(item))
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
		helper.JsonResponse(w, http.StatusInternalServerError, "Error updating item")
		return
	}

	helper.JsonResponse(w, http.StatusNoContent, nil)
}

func (h *Handler) decodeItemRequest(w http.ResponseWriter, r *http.Request) (*ItemRequest, error) {
	itemRequest := ItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(&itemRequest); err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
		return nil, err
	}

	return &itemRequest, nil
}

func (h *Handler) getItem(w http.ResponseWriter, r *http.Request) (*Item, error) {
	itemId, err := helper.ConvertStringIdToUint(chi.URLParam(r, "itemId"))
	if err != nil {
		helper.ErrorResponse(w, http.StatusBadRequest, "Invalid id format")
		return nil, err
	}

	item, err := h.repository.Get(itemId)
	if err == ErrItemNotFound {
		helper.ErrorResponse(w, http.StatusNotFound, "Not found")
		return nil, err
	} else if err != nil {
		log.Printf("Error: %s", err.Error())
		helper.ErrorResponse(w, http.StatusInternalServerError, "Error fetching item")
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
