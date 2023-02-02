package handlers_test

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/osemisan/osemisan-client/pkg/handlers"
	"github.com/osemisan/osemisan-client/testutil"
)

func MockAuthorizationServer(t *testing.T) (*httptest.Server, error) {
	r := chi.NewRouter()

	r.Get("/authorize", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	r.Post("/token", func(w http.ResponseWriter, r *http.Request) {
		tokenRes := handlers.TokenResponse{
			AccessToken: testutil.BuildScopedJwt(t, testutil.Scopes{
				Abura:  true,
				Minmin: true,
			}),
			TokenType: "Bearer",
		}
		bytes, _ := json.Marshal(tokenRes)
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	})

	l, err := net.Listen("tcp", ":9001")
	if err != nil {
		return nil, err
	}
	ts := httptest.Server{
		Listener: l,
		Config:   &http.Server{Handler: r},
	}
	return &ts, nil
}
