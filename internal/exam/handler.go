package exam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateExam godoc
// @Summary Create a new exam with questions list
// @Description Creates a new exam instance with an associated questions list in a single transaction. Accepts exam details and a list of questions to be linked to the exam.
// @Tags Exams
// @Accept json
// @Produce json
// @Param exam body NewExam true "Exam instance and questions list details"
// @Success 201 {object} map[string]interface{} "Successfully created exam"
// @Failure 400 {string} string "Invalid request body or malformed JSON"
// @Failure 500 {string} string "Internal server error during exam creation"
// @Router /exams [post]
func CreateExam(w http.ResponseWriter, r *http.Request) {
	var newExam NewExam

	err := json.NewDecoder(r.Body).Decode(&newExam)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		println(err.Error())
		return
	}

	response, err := createExamPayload(newExam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating exam: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
