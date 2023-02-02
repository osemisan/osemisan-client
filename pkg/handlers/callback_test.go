package handlers_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/osemisan/osemisan-client/pkg/handlers"
	"github.com/osemisan/osemisan-client/testutil"
)

func TestCallbackHandler(t *testing.T) {
	dummyState := "dummystate"

	tests := []struct {
		name                    string
		query                   url.Values
		wantTextInResHTML       string
		wantAccessTokenInCookie string
		wantStatusCode          int
	}{
		{
			name: "クエリパラメターにerrorをつけると、エラーページが返ってくる",
			query: url.Values{
				"state": {dummyState},
				"error": {"ERROR MESSAGE"},
			},
			wantTextInResHTML:       "ERROR MESSAGE",
			wantAccessTokenInCookie: "",
			wantStatusCode:          http.StatusOK,
		},
		{
			name: "不正なstate値を渡すと、エラーページが返ってくる",
			query: url.Values{
				"state": {"invaliddummystate"},
			},
			wantTextInResHTML:       "stateがマッチしません",
			wantAccessTokenInCookie: "",
			wantStatusCode:          http.StatusOK,
		},
		{
			name: "リクエストが正常に処理された場合、クッキーにアクセストークンが格納されている",
			query: url.Values{
				"state": {dummyState},
			},
			wantTextInResHTML: "アクセストークン",
			wantAccessTokenInCookie: testutil.BuildScopedJwt(t, testutil.Scopes{
				Abura:  true,
				Minmin: true,
			}),
			wantStatusCode: http.StatusOK,
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
			s := httptest.NewServer(http.HandlerFunc(handlers.CallbackHandler))
			u, err := url.Parse(s.URL)
			if err != nil {
				t.Error("Failed to parse URL", err)
			}
			u.RawQuery = tt.query.Encode()

			req, err := http.NewRequest(http.MethodGet, u.String(), nil)
			if err != nil {
				t.Error("Failed to create new request", err)
				return
			}

			// すでに authorize エンドポイントを叩いている想定なので、ステートをクッキーにセットしておく
			req.AddCookie(&http.Cookie{
				Name:  "state",
				Value: dummyState,
			})

			res, err := c.Do(req)
			if err != nil {
				t.Error("Failed to request", err)
				return
			}
			defer res.Body.Close()

			if tt.wantStatusCode != res.StatusCode {
				t.Errorf("Unexpected statuc code, expected: %d, actual: %d", tt.wantStatusCode, res.StatusCode)
				return
			}

			body, _ := ioutil.ReadAll(res.Body)
			buf := bytes.NewBuffer(body)
			html := buf.String()

			if !strings.Contains(html, tt.wantTextInResHTML) {
				t.Errorf(`Response HTML does not contain "%s"`, tt.wantTextInResHTML)
				return
			}

			if tt.wantAccessTokenInCookie != "" && tt.wantAccessTokenInCookie != res.Cookies()[0].Value {
				t.Errorf("Unexpected tokein in Cookie, expected %s, actual: %s", tt.wantAccessTokenInCookie, res.Cookies()[0].Value)
				return
			}
		})
	}
}
