package answer

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes wires the answer domain dependencies and registers its
// routes. The given router is expected to already carry the authentication
// middleware.
func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	handler := NewHandler(NewService(NewRepository(db)))

	router.HandleFunc("/answers/performance", handler.GetSubjectPerfomanceHandler).Methods("GET")
	router.HandleFunc("/answers/overall-performance", handler.GetUserOverallPerfomanceHandler).Methods("GET")
}
