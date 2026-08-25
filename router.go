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
	"flashquest/internal/questionoption"
	"flashquest/internal/questionset"
	"flashquest/internal/source"
	"flashquest/internal/submission"
	"flashquest/internal/user"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"

	_ "flashquest/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewServer(db *gorm.DB) *http.Server {
	r := mux.NewRouter()
	registerPaths(r, db)

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

func registerPaths(r *mux.Router, db *gorm.DB) {
	// Auth Requests
	authRepo := auth.NewRepositoryWithDB(db)
	authService := auth.NewService(authRepo, os.Getenv("JWT_PRIVATE_KEY"))
	r.HandleFunc("/api/auth/register", auth.RegisterHandler(authService)).Methods("POST")
	r.HandleFunc("/api/auth/login", auth.LoginHandler(authService)).Methods("POST")

	// Question Requests
	questionHandler := question.NewHandler(question.NewRepository(db))
	r.HandleFunc("/questions", questionHandler.CreateQuestion).Methods("POST") //Updated
	r.HandleFunc("/questions", questionHandler.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/filters", questionHandler.GetQuestionFiltersHandler).Methods("GET")
	r.HandleFunc("/questions/by-ids", questionHandler.GetQuestionsByArray).Methods("POST")
	r.HandleFunc("/questions/{id}", questionHandler.GetQuestion).Methods("GET")
	r.HandleFunc("/questions/{id}", questionHandler.DeleteQuestion).Methods("DELETE")

	// QuestionOption Requests
	questionOptionHandler := questionoption.NewHandler(questionoption.NewRepository(db))
	r.HandleFunc("/questions/{id}/question-options", questionOptionHandler.PostQuestionOptions).Methods("POST")
	r.HandleFunc("/questions/{id}/question-options", questionOptionHandler.GetQuestionOptions).Methods("GET")
	r.HandleFunc("/question-options/by-ids", questionOptionHandler.GetQuestionOptionsByIDArray).Methods("POST")

	//Question Set Requests
	questionSetHandler := questionset.NewHandler(questionset.NewRepository(db))
	r.HandleFunc("/question-sets", questionSetHandler.CreateQuestionSet).Methods("POST")
	r.HandleFunc("/question-sets", questionSetHandler.GetLists).Methods("GET")
	r.HandleFunc("/question-sets/{id}", questionSetHandler.GetQuestionSet).Methods("GET")
	r.HandleFunc("/question-sets/{id}/questions", questionSetHandler.GetQuestionsFromSet).Methods("GET")

	sourceHandler := source.NewHandler(source.NewRepository(db))
	r.HandleFunc("/sources", sourceHandler.CreateSource).Methods("POST")

	//AI requests
	aiHandler := ai.NewHandler(ai.NewService(question.NewRepository(db), questionoption.NewRepository(db), questionset.NewRepository(db)))
	r.HandleFunc("/ai/gen-question", aiHandler.PostAIGenQuestion).Methods("POST")
	r.HandleFunc("/ai/gen-questionset", aiHandler.PostAIGenQuestionSet).Methods("POST")

	examHandler := exam.NewHandler(exam.NewRepository(db), question.NewRepository(db), questionoption.NewRepository(db), questionset.NewRepository(db))
	r.HandleFunc("/exam", examHandler.CreateExam).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	protectedRouter := r.PathPrefix("/").Subrouter()
	protectedRouter.Use(auth.RequireAuth)

	submissionHandler := submission.NewHandler(submission.NewRepository(db))
	protectedRouter.HandleFunc("/submissions", submissionHandler.CreateSubmission).Methods("POST")
	protectedRouter.HandleFunc("/submissions", submissionHandler.GetUserSubmissions).Methods("GET")
	protectedRouter.HandleFunc("/submissions/{id}", submissionHandler.GetSubmission).Methods("GET")
	userHandler := user.NewHandler(user.NewService(user.NewRepository(db)))
	protectedRouter.HandleFunc("/users/me", userHandler.GetCurrentUser).Methods("GET")
	answerHandler := answer.NewHandler(answer.NewRepository(db))
	protectedRouter.HandleFunc("/answers/performance", answerHandler.GetSubjectPerfomanceHandler).Methods("GET")
	protectedRouter.HandleFunc("/answers/overall-performance", answerHandler.GetUserOverallPerfomanceHandler).Methods("GET")
}
