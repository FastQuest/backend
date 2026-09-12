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

	var m mux.RouteMatch
	if !r.Match(httptest.NewRequest(http.MethodGet, "/users/me", nil), &m) {
		t.Fatal("expected /users/me route to be registered")
	}
}
