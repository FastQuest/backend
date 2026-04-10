package ai

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func PostAIGenQuestion(w http.ResponseWriter, r *http.Request) {
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
	addAIQuestion(genQuestion(test.Text))
}

func PostAIGenQuestionSet(w http.ResponseWriter, r *http.Request) {
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

	errAddQS := addAIQuestionSet(genQuestionSet(test.Text))
	if errAddQS != nil {
		http.Error(w, "Failed to generate question set", http.StatusInternalServerError)
	}
}
