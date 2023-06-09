package helper

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func ConvertStringIdToUint(idParam string) (uint, error) {
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func JsonResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
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

func ErrorResponse(w http.ResponseWriter, statusCode int, errorMessage string) {
	JsonResponse(w, statusCode, map[string]string{"error": errorMessage})
}
