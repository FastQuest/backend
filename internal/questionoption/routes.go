package questionoption

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, service Service) {
	handler := NewHandler(service)

	router.HandleFunc("/questions/{id}/question-options", handler.PostQuestionOptions).Methods("POST")
	router.HandleFunc("/questions/{id}/question-options", handler.GetQuestionOptions).Methods("GET")
	router.HandleFunc("/question-options/by-ids", handler.GetQuestionOptionsByIDArray).Methods("POST")
}
