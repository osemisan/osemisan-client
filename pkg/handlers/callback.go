package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/httplog"
	"github.com/osemisan/osemisan-client/pkg/client"
	"github.com/osemisan/osemisan-client/pkg/endpoints"
	"github.com/osemisan/osemisan-client/pkg/scope"
	"github.com/osemisan/osemisan-client/pkg/templates"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	r.ParseForm()

	e := r.FormValue("error")

	if e != "" {
		err := templates.Render("error", w, map[string]string{
			"error": e,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	state := r.FormValue("state")

	stateCookie, err := r.Cookie("state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if state == stateCookie.Value {
		oplog.Info().Msg("State matched")
	} else {
		err := templates.Render("error", w, map[string]string{
			"error": "state がマッチしません",
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	code := r.FormValue("code")

	f := url.Values{}
	f.Add("grant_type", "authorization_code")
	f.Add("code", code)
	f.Add("redirect_uri", client.C.URIs[0])

	c := new(http.Client)

	b := bytes.NewBuffer([]byte(f.Encode()))
	req, err := http.NewRequest(http.MethodPost, endpoints.TokenEndpoint, b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.SetBasicAuth(client.C.Id, client.C.Secret)

	tokRes, err := c.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer tokRes.Body.Close()

	oplog.Info().Msgf("Requesting access token for code %s", code)

	if tokRes.StatusCode == http.StatusOK {
		tok := new(TokenResponse)

		body, err := io.ReadAll(tokRes.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(body, &tok); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		c := &http.Cookie{
			Name:  "token",
			Value: tok.AccessToken,
		}
		http.SetCookie(w, c)

		scope, err := scope.Store.Get(tok.AccessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := templates.Render("index", w, map[string]string{
			"accessToken": tok.AccessToken,
			"scope":       scope,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		err := templates.Render("error", w, map[string]string{
			"error": fmt.Sprintf("アクセストークンのフェッチに失敗しました。サーバーレスポンス: %d", tokRes.StatusCode),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
