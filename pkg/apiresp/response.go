package apiresp

import (
	"encoding/json"
	"net/http"
)

const fallbackErrorResponse = "{\"error\":{\"code\":\"INTERNAL_SERVER_ERROR\",\"message\":\"internal server error\"}}\n"

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, code, message string) {
	resp := errorEnvelope{}
	resp.Error.Code = code
	resp.Error.Message = message

	writeJSON(w, status, resp)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	writeJSON(w, status, payload)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		writeFallbackInternalError(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}

func writeFallbackInternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte(fallbackErrorResponse))
}
