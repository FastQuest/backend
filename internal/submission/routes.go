package submission

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, service Service) {
	handler := NewHandler(service)

	router.HandleFunc("/submissions", handler.CreateSubmission).Methods("POST")
	router.HandleFunc("/submissions", handler.GetUserSubmissions).Methods("GET")
	router.HandleFunc("/submissions/{id}", handler.GetSubmission).Methods("GET")
}
