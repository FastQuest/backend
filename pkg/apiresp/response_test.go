package apiresp_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"flashquest/pkg/apiresp"
)

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()

	apiresp.WriteError(rr, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Email ou senha inválidos")

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	expected := "{\"error\":{\"code\":\"INVALID_CREDENTIALS\",\"message\":\"Email ou senha inválidos\"}}\n"
	if rr.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rr.Body.String())
	}
}

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	payload := map[string]any{"ok": true}
	apiresp.WriteJSON(rr, http.StatusCreated, payload)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	expected := "{\"ok\":true}\n"
	if rr.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rr.Body.String())
	}
}
