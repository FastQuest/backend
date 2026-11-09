package main

import (
	"net/http"
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
	auth.RegisterRoutes(r, db)

	questionService := question.NewService(question.NewRepository(db))
	questionOptionService := questionoption.NewService(questionoption.NewRepository(db), questionService)

	// Question Requests
	question.RegisterRoutes(r, db)

	// QuestionOption Requests
	questionoption.RegisterRoutes(r, questionOptionService)

	//Question Set Requests
	questionSetHandler := questionset.NewHandler(questionset.NewRepository(db))
	r.HandleFunc("/question-sets", questionSetHandler.CreateQuestionSet).Methods("POST")
	r.HandleFunc("/question-sets", questionSetHandler.GetLists).Methods("GET")
	r.HandleFunc("/question-sets/{id}", questionSetHandler.GetQuestionSet).Methods("GET")
	r.HandleFunc("/question-sets/{id}/questions", questionSetHandler.GetQuestionsFromSet).Methods("GET")

	sourceHandler := source.NewHandler(source.NewRepository(db))
	r.HandleFunc("/sources", sourceHandler.CreateSource).Methods("POST")

	//AI requests
	aiHandler := ai.NewHandler(ai.NewService(questionService, questionOptionService, questionset.NewRepository(db)))
	r.HandleFunc("/ai/gen-question", aiHandler.PostAIGenQuestion).Methods("POST")
	r.HandleFunc("/ai/gen-questionset", aiHandler.PostAIGenQuestionSet).Methods("POST")

	examHandler := exam.NewHandler(exam.NewRepository(db), questionService, questionOptionService, questionset.NewRepository(db))
	r.HandleFunc("/exam", examHandler.CreateExam).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	protectedRouter := r.PathPrefix("/").Subrouter()
	protectedRouter.Use(auth.RequireAuth)

	answerService := answer.NewService(answer.NewRepository(db))
	submissionService := submission.NewService(submission.NewRepository(db), answerService, questionOptionService)
	submissionHandler := submission.NewHandler(submissionService)
	protectedRouter.HandleFunc("/submissions", submissionHandler.CreateSubmission).Methods("POST")
	protectedRouter.HandleFunc("/submissions", submissionHandler.GetUserSubmissions).Methods("GET")
	protectedRouter.HandleFunc("/submissions/{id}", submissionHandler.GetSubmission).Methods("GET")
	user.RegisterRoutes(protectedRouter, db)
	answer.RegisterRoutes(protectedRouter, answerService)
}
