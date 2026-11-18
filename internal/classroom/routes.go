package classroom

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router, service Service) {
	handler := NewHandler(service)
	router.HandleFunc("/classrooms", handler.CreateClassroom).Methods("POST")
	router.HandleFunc("/classrooms", handler.ListClassrooms).Methods("GET")
	router.HandleFunc("/classrooms/{id}/status", handler.UpdateClassroomStatus).Methods("PUT")
	router.HandleFunc("/classrooms/{id}/participants", handler.AddParticipant).Methods("POST")
	router.HandleFunc("/classrooms/{id}/participants", handler.ListParticipants).Methods("GET")
	router.HandleFunc("/classrooms/{id}/activities", handler.CreateActivity).Methods("POST")
	router.HandleFunc("/classrooms/{id}/activities", handler.ListActivities).Methods("GET")
}
