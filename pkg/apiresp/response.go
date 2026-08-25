package apiresp

import (
	"encoding/json"
	"log"
	"net/http"
)

const fallbackErrorResponse = "{\"error\":{\"code\":\"INTERNAL_SERVER_ERROR\",\"message\":\"internal server error\"}}\n"

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

type paginatedEnvelope struct {
	Data any        `json:"data"`
	Meta Pagination `json:"meta"`
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

// WritePaginatedJSON writes a list response using the standard envelope:
// the items under "data" and the pagination info under "meta".
func WritePaginatedJSON(w http.ResponseWriter, status int, data any, meta Pagination) {
	writeJSON(w, status, paginatedEnvelope{Data: data, Meta: meta})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("apiresp: failed to marshal payload of type %T: %v", payload, err)
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
