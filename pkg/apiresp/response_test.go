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

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
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

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}
}

func TestWriteJSONEncodeFailureFallback(t *testing.T) {
	rr := httptest.NewRecorder()

	apiresp.WriteJSON(rr, http.StatusCreated, map[string]any{"invalid": make(chan int)})

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	expected := "{\"error\":{\"code\":\"INTERNAL_SERVER_ERROR\",\"message\":\"internal server error\"}}\n"
	if rr.Body.String() != expected {
		t.Fatalf("expected deterministic fallback body %q, got %q", expected, rr.Body.String())
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}
}
