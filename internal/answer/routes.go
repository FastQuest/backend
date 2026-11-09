package answer

import (
	"github.com/gorilla/mux"
)

// RegisterRoutes wires the answer domain dependencies and registers its
// routes. The given router is expected to already carry the authentication
// middleware.
func RegisterRoutes(router *mux.Router, service Service) {
	handler := NewHandler(service)

	router.HandleFunc("/answers/performance", handler.GetSubjectPerfomanceHandler).Methods("GET")
	router.HandleFunc("/answers/overall-performance", handler.GetUserOverallPerfomanceHandler).Methods("GET")
}
