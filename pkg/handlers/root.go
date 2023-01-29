package handlers

import (
	"net/http"

	"github.com/osemisan/osemisan-client/pkg/templates"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.Render("index", w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
