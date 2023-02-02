package handlers_test

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/osemisan/osemisan-client/pkg/handlers"
)

func MockResourceServer() (*httptest.Server, error) {
	r := chi.NewRouter()

	r.Get("/resources", func(w http.ResponseWriter, r *http.Request) {
		resources := []handlers.SemiResource{
			{Name: "アブラゼミ", Length: "5cm"},
			{Name: "クマゼミ", Length: "7cm"},
		}
		bytes, _ := json.Marshal(resources)
		w.Header().Set("Content-Type", "application/json")
		w.Write(bytes)
	})

	l, err := net.Listen("tcp", ":9002")
	if err != nil {
		return nil, err
	}
	ts := &httptest.Server{
		Listener: l,
		Config:   &http.Server{Handler: r},
	}
	return ts, nil
}
