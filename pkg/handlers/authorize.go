package handlers

import (
	"net/http"
	"net/url"

	"github.com/osemisan/osemisan-client/pkg/client"
	"github.com/osemisan/osemisan-client/pkg/endpoints"
	"github.com/osemisan/osemisan-client/pkg/random"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	// MaxAgeを負数にして期限切れにして実質削除する
	c := &http.Cookie{
		Name:   "token",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// ランダムなステートを生成してクッキーに入れとく
	state := random.GenStr(32)
	stateC := http.Cookie{
		Name:  "state",
		Value: state,
	}
	http.SetCookie(w, &stateC)

	u, err := url.Parse(endpoints.AuthorizationEndpoint)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	q := u.Query()
	q.Set("response_type", "code")
	q.Set("scope", client.C.Scope)
	q.Set("client_id", client.C.Id)
	q.Set("redirect_uri", client.C.URIs[0])
	q.Set("state", state)
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusPermanentRedirect)
}
