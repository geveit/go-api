package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func WriteJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		// Error encoding JSON response
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func WriteError(w http.ResponseWriter, statusCode int, errorMessage string) {
	WriteJSON(w, statusCode, map[string]string{"error": errorMessage})
}

func GetIdFromParam(r *http.Request, param string) (uint, error) {
	id, err := strconv.ParseUint(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

type ApiServer struct {
	port   string
	Router chi.Router
}

func NewServer(port string) *ApiServer {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	return &ApiServer{port: port, Router: router}
}

func (s *ApiServer) Run() {
	log.Printf("Server running on port %s", s.port)

	log.Fatal(http.ListenAndServe(s.port, s.Router))
}
