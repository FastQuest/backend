package questionset

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, service Service) {
	handler := NewHandler(service)

	router.HandleFunc("/question-sets", handler.CreateQuestionSet).Methods("POST")
	router.HandleFunc("/question-sets", handler.GetLists).Methods("GET")
	router.HandleFunc("/question-sets/{id}", handler.GetQuestionSet).Methods("GET")
	router.HandleFunc("/question-sets/{id}/questions", handler.GetQuestionsFromSet).Methods("GET")
}
