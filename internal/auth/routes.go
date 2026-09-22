package auth

import (
	"os"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes wires the auth domain dependencies and registers its routes.
func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	handler := NewHandler(NewService(NewRepositoryWithDB(db), os.Getenv("JWT_PRIVATE_KEY")))

	router.HandleFunc("/api/auth/register", handler.Register).Methods("POST")
	router.HandleFunc("/api/auth/login", handler.Login).Methods("POST")
}
