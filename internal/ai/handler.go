package ai

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) PostAIGenQuestion(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var test TestBody
	errConvert := json.Unmarshal(body, &test)
	if errConvert != nil {
		http.Error(w, "Invalid body", http.StatusInternalServerError)
		return
	}

	log.Println("Successful POST")
	h.service.addAIQuestion(genQuestion(test.Text))
}

// GetQuestions godoc
// @Summary Generate question set by AI
// @Description Generate question set by AI
// @Tags AI
// @Accept json
// @Produce json
// @Param question body TestBody true "Text to generate a set of questions"
// @Success 201
// @Failure 500 {string} string "Internal server error"
// @Router /ai/gen-questionset [post]
func (h *Handler) PostAIGenQuestionSet(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var test TestBody
	errConvert := json.Unmarshal(body, &test)
	if errConvert != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	errAddQS := h.service.addAIQuestionSet(genQuestionSet(test.Text))
	if errAddQS != nil {
		http.Error(w, "Failed to generate question set", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}
