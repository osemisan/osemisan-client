package handlers

import (
	"errors"
	"net/http"

	"github.com/osemisan/osemisan-client/pkg/scope"
	"github.com/osemisan/osemisan-client/pkg/templates"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	tok := ""
	cookie, err := r.Cookie("token")
	if err != nil {
		if !errors.Is(err, http.ErrNoCookie) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
	  tok = cookie.Value
	}
	var s string
	if tok != "" {
		s, err = scope.Store.Get(tok)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	err = templates.Render("index", w, map[string]string {
		"accessToken": tok,
		"scope": s,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
