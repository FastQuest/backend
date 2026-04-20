package exam

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
