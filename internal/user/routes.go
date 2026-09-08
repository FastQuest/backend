package user

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes wires the user domain dependencies and registers its routes.
// The given router is expected to already carry the authentication middleware.
func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	handler := NewHandler(NewService(NewRepository(db)))

	router.HandleFunc("/users/me", handler.GetCurrentUser).Methods("GET")
}
