package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/osemisan/osemisan-client/pkg/handlers"
)

func main() {
	l := httplog.NewLogger("osemisan-client", httplog.Options{
		JSON: true,
	})
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(l))

	r.Get("/", handlers.RootHandler)
	r.Get("/authorize", handlers.AuthorizeHandler)
	r.Get("/callback", handlers.CallbackHandler)
	r.Get("/fetch_resource", handlers.FetchResourceHandler)

	http.ListenAndServe(":9000", r)
}
