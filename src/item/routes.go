package item

import "github.com/go-chi/chi"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/items", func(r chi.Router) {
		r.Route("/{itemId}", func(r chi.Router) {
			r.Get("/", h.GetItem)
			r.Post("/", h.UpdateItem)
			r.Delete("/", h.DeleteItem)
		})

		r.Get("/", h.GetAllItems)
		r.Post("/", h.CreateItem)
	})
}
