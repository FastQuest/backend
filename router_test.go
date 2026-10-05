package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestRoutesWiring(t *testing.T) {
	r := mux.NewRouter()
	registerPaths(r, nil)

	cases := []struct {
		name   string
		method string
		path   string
		body   string
		want   int
	}{
		{"register invalid payload", "POST", "/api/auth/register", "not-json", 400},
		{"register invalid email", "POST", "/api/auth/register", `{"name":"a","email":"nope","password":"password123"}`, 400},
		{"register bad domain", "POST", "/api/auth/register", `{"name":"a","email":"a@gmail.com","password":"password123"}`, 422},
		{"login invalid payload", "POST", "/api/auth/login", "not-json", 400},
		{"login missing password", "POST", "/api/auth/login", `{"email":"a@ceub.edu.br","password":""}`, 400},
		{"users me unauthenticated", "GET", "/users/me", "", 401},
		{"subject performance unauthenticated", "GET", "/answers/performance", "", 401},
		{"overall performance unauthenticated", "GET", "/answers/overall-performance", "", 401},
		{"create question invalid payload", "POST", "/questions", "not-json", 400},
		{"create question rejects array payload", "POST", "/questions", `[{"statement":"a","user_id":1}]`, 400},
		{"create question missing fields", "POST", "/questions", `{"statement":"a"}`, 400},
		{"batch questions invalid payload", "POST", "/questions/batch", "not-json", 400},
		{"batch questions rejects object payload", "POST", "/questions/batch", `{"statement":"a","user_id":1}`, 400},
		{"batch questions empty", "POST", "/questions/batch", `[]`, 400},
		{"batch questions missing fields", "POST", "/questions/batch", `[{"statement":"","user_id":1}]`, 400},
		{"questions by ids empty", "POST", "/questions/by-ids", `{"ids":[]}`, 400},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)
			if rec.Code != tc.want {
				t.Fatalf("%s %s: expected %d, got %d body=%s", tc.method, tc.path, tc.want, rec.Code, rec.Body.String())
			}
		})
	}

	for _, path := range []string{"/users/me", "/answers/performance", "/answers/overall-performance"} {
		var m mux.RouteMatch
		if !r.Match(httptest.NewRequest(http.MethodGet, path, nil), &m) {
			t.Fatalf("expected %s route to be registered", path)
		}
	}

	questionRoutes := []struct{ method, path string }{
		{http.MethodPost, "/questions"},
		{http.MethodPost, "/questions/batch"},
		{http.MethodGet, "/questions"},
		{http.MethodGet, "/questions/filters"},
		{http.MethodPost, "/questions/by-ids"},
		{http.MethodGet, "/questions/1"},
		{http.MethodDelete, "/questions/1"},
	}
	for _, route := range questionRoutes {
		var m mux.RouteMatch
		if !r.Match(httptest.NewRequest(route.method, route.path, nil), &m) {
			t.Fatalf("expected %s %s route to be registered", route.method, route.path)
		}
	}
}
