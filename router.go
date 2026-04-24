package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"flashquest/internal/ai"
	"flashquest/internal/answer"
	"flashquest/internal/auth"
	"flashquest/internal/exam"
	"flashquest/internal/question"
	"flashquest/internal/questionset"
	"flashquest/internal/source"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "flashquest/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewServer() *http.Server {
	r := mux.NewRouter()
	registerPaths(r)

	c := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		AllowedOrigins:   []string{"https://fastquest.vercel.app"},

		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
	})

	handler := c.Handler(r)

	return &http.Server{
		Handler:      handler,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

func registerPaths(r *mux.Router) {
	// Auth Requests
	authRepo := auth.NewRepository()
	authService := auth.NewService(authRepo, os.Getenv("JWT_PRIVATE_KEY"))
	r.HandleFunc("/api/auth/register", auth.RegisterHandler(authService)).Methods("POST")
	r.HandleFunc("/api/auth/login", auth.LoginHandler(authService)).Methods("POST")

	// Question Requests
	r.HandleFunc("/questions", question.CreateQuestion).Methods("POST") //Updated
	r.HandleFunc("/questions", question.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/by-ids", question.GetQuestionsByArray).Methods("POST")
	r.HandleFunc("/questions/{id}", question.GetQuestion).Methods("GET")
	r.HandleFunc("/questions/{id}", question.DeleteQuestion).Methods("DELETE")

	// Answer Requests
	r.HandleFunc("/questions/{id}/answers", answer.PostAnswers).Methods("POST")
	r.HandleFunc("/questions/{id}/answers", answer.GetAnswers).Methods("GET")
	r.HandleFunc("/answers/by-ids", answer.GetAnswersByIDArray).Methods("POST")

	//Question Set Requests
	r.HandleFunc("/question-sets", questionset.CreateQuestionSet).Methods("POST")
	r.HandleFunc("/question-sets", questionset.GetLists).Methods("GET")
	r.HandleFunc("/question-sets/{id}", questionset.GetQuestionSet).Methods("GET")
	r.HandleFunc("/question-sets/{id}/questions", questionset.GetQuestionsFromSet).Methods("GET")

	r.HandleFunc("/sources", source.CreateSource).Methods("POST")

	//AI requests
	r.HandleFunc("/ai/gen-question", ai.PostAIGenQuestion).Methods("POST")
	r.HandleFunc("/ai/gen-questionset", ai.PostAIGenQuestionSet).Methods("POST")

	r.HandleFunc("/exam", exam.CreateExam).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
