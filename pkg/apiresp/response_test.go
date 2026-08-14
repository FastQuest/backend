package apiresp_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestWritePaginatedJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	items := []map[string]any{{"id": 1}, {"id": 2}}
	apiresp.WritePaginatedJSON(rr, http.StatusOK, items, apiresp.NewPagination(25, 2, 10))

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	expected := "{\"data\":[{\"id\":1},{\"id\":2}],\"meta\":{\"total\":25,\"per_page\":10,\"current_page\":2,\"last_page\":3}}\n"
	if rr.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rr.Body.String())
	}

	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}
}

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name             string
		total            int64
		currentPage      int
		perPage          int
		expectedLastPage int
	}{
		{name: "exact division", total: 30, currentPage: 1, perPage: 10, expectedLastPage: 3},
		{name: "partial last page", total: 31, currentPage: 1, perPage: 10, expectedLastPage: 4},
		{name: "no results", total: 0, currentPage: 1, perPage: 10, expectedLastPage: 0},
		{name: "zero per page does not divide by zero", total: 10, currentPage: 1, perPage: 0, expectedLastPage: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := apiresp.NewPagination(tt.total, tt.currentPage, tt.perPage)

			if got.LastPage != tt.expectedLastPage {
				t.Fatalf("expected last_page %d, got %d", tt.expectedLastPage, got.LastPage)
			}
			if got.Total != int(tt.total) {
				t.Fatalf("expected total %d, got %d", tt.total, got.Total)
			}
			if got.PerPage != tt.perPage {
				t.Fatalf("expected per_page %d, got %d", tt.perPage, got.PerPage)
			}
			if got.CurrentPage != tt.currentPage {
				t.Fatalf("expected current_page %d, got %d", tt.currentPage, got.CurrentPage)
			}
		})
	}
}

func TestWriteJSONEncodeFailureFallback(t *testing.T) {
	rr := httptest.NewRecorder()

	var logs bytes.Buffer
	log.SetOutput(&logs)
	defer log.SetOutput(os.Stderr)

	apiresp.WriteJSON(rr, http.StatusCreated, map[string]any{"invalid": make(chan int)})

	if !strings.Contains(logs.String(), "failed to marshal payload") {
		t.Fatalf("expected marshal failure to be logged, got %q", logs.String())
	}

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
