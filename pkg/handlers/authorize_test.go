package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osemisan/osemisan-client/pkg/client"
	"github.com/osemisan/osemisan-client/pkg/handlers"
)

func TestAuthorizeHandler(t *testing.T) {
	tests := []struct {
		name           string
		withTok        bool
		wantStatusCode int
	}{
		{
			"クッキーのトークンが格納された状態でリクエストしたらクッキーのトークンが削除される",
			true,
			http.StatusFound,
		},
		{
			"クッキーにトークンがないときにリクエストしたら当然クッキーにトークンはない",
			false,
			http.StatusFound,
		},
	}
	c := new(http.Client)

	ts, err := MockAuthorizationServer(t)
	if err != nil {
		t.Error("Failed to create mock server", err)
		return
	}
	ts.Start()
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(handlers.AuthorizeHandler))

			req, err := http.NewRequest(http.MethodGet, s.URL, nil)
			if err != nil {
				t.Error("Failed to create new request", err)
			}

			if tt.withTok {
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: "dummytoken",
				})
			}

			res, err := c.Do(req)
			if err != nil {
				t.Error("Failed to request", err)
				return
			}

			cs := res.Request.Response.Cookies()

			// レスポンスのクッキーからはトークンが消えている
			tokC := cs[0]
			if tokC.MaxAge != -1 {
				t.Error(`"token" isn't expired`)
				return
			}
			if tokC.Value != "" {
				t.Error(`"token" isn't empty`)
				return
			}

			// レスポンスのクッキーに32ケタのステートが格納されている
			stateC := cs[1]
			if len(stateC.Value) != 32 {
				t.Error(`"state" is invalid`)
				return
			}

			// リダイレクト時のURLにちゃんとクエリパラメータが指定されているかどうか
			q := res.Request.URL.Query()

			gotId := q.Get("client_id")
			if gotId != client.C.Id {
				t.Errorf("Unexpected client ID, expected: %s, actual: %s", client.C.Id, gotId)
				return
			}

			gotURI := q.Get("redirect_uri")
			if gotURI != client.C.URIs[0] {
				t.Errorf("Unexpected redirect URI, expected: %s, actual: %s", client.C.URIs[0], gotURI)
				return
			}

			gotType := q.Get("response_type")
			if gotType != "code" {
				t.Errorf("Unexpected response type, expected: code, actual: %s", gotType)
				return
			}

			gotScope := q.Get("scope")
			if gotScope != client.C.Scope {
				t.Errorf("Unexpected scope, expected: %s, actual: %s", client.C.Scope, gotScope)
				return
			}

			gotState := q.Get("state")
			if len(gotState) != 32 {
				t.Errorf("Invalid state in URI query, %s", gotState)
				return
			}

			if res.Request.Response.StatusCode != tt.wantStatusCode {
				t.Errorf("Unexpected status code, expected: %d, actual: %d", tt.wantStatusCode, res.StatusCode)
				return
			}
		})
	}
}
