package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/httplog"
	"github.com/osemisan/osemisan-client/pkg/endpoints"
	"github.com/osemisan/osemisan-client/pkg/templates"
)

type SemiResource struct {
	Name   string `json:"name"`
	Length string `json:"length"`
}

func FetchResourceHandler(w http.ResponseWriter, r *http.Request) {
	oplog := httplog.LogEntry(r.Context())

	cookie, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			if err := templates.Render("error", w, map[string]string{
				"message": "Missing access token",
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := cookie.Value

	oplog.Info().Msgf("Missing request with access token %s", token)

	c := new(http.Client)

	req, err := http.NewRequest(http.MethodGet, endpoints.ResoucesEndpoint, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if res.StatusCode == http.StatusOK {
		r := new([]SemiResource)

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(body, &r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := templates.Render("data", w, map[string]any{
			"resource": &r,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// 格納されているアクセストークンをクッキーからクリアする
		http.SetCookie(w, &http.Cookie{
			Name:   "token",
			Value:  "",
			MaxAge: -1,
		})
		if err := templates.Render("data", w, map[string]string{
			"message": fmt.Sprintf("リソースサーバーが %d を返しました", res.StatusCode),
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
