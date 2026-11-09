package question

import (
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// RegisterRoutes wires the question domain dependencies and registers its
// routes.
func RegisterRoutes(router *mux.Router, db *gorm.DB) {
	handler := NewHandler(NewService(NewRepository(db)))

	router.HandleFunc("/questions", handler.CreateQuestion).Methods("POST")
	router.HandleFunc("/questions/batch", handler.CreateQuestionsBatch).Methods("POST")
	router.HandleFunc("/questions", handler.GetQuestions).Methods("GET")
	router.HandleFunc("/questions/filters", handler.GetQuestionFiltersHandler).Methods("GET")
	router.HandleFunc("/questions/by-ids", handler.GetQuestionsByArray).Methods("POST")
	router.HandleFunc("/questions/{id}", handler.GetQuestion).Methods("GET")
	router.HandleFunc("/questions/{id}", handler.DeleteQuestion).Methods("DELETE")
}
