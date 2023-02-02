package handlers_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/osemisan/osemisan-client/pkg/handlers"
	"github.com/osemisan/osemisan-client/testutil"
)

func TestFetchResourceHandler(t *testing.T) {
	tests := []struct {
		name              string
		reqWithToken      bool
		wantTextInResHTML string
		wantStatusCode    int
	}{
		{
			name:              "Cookieにトークンが格納されていない状態でリクエストをすると、エラーページを返す",
			reqWithToken:      false,
			wantTextInResHTML: "Missing access token",
			wantStatusCode:    http.StatusOK,
		},
		{
			name:              "Cookieにトークンが格納されている状態でリクエストすると、リソースを表示する",
			reqWithToken:      true,
			wantTextInResHTML: "アブラ",
			wantStatusCode:    http.StatusOK,
		},
	}

	c := new(http.Client)

	ts, err := MockResourceServer()
	if err != nil {
		t.Error("Failed to create mock server", err)
		return
	}
	ts.Start()
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(handlers.FetchResourceHandler))

			req, err := http.NewRequest(http.MethodGet, s.URL, nil)
			if err != nil {
				t.Error("Failed to create new request", err)
				return
			}

			if tt.reqWithToken {
				req.AddCookie(&http.Cookie{
					Name: "token",
					Value: testutil.BuildScopedJwt(t, testutil.Scopes{
						Abura: true,
						Kuma:  true,
					}),
				})
			}

			res, err := c.Do(req)
			if err != nil {
				t.Error("Failed to request", err)
				return
			}
			defer res.Body.Close()

			body, _ := ioutil.ReadAll(res.Body)
			buf := bytes.NewBuffer(body)
			html := buf.String()

			if tt.wantStatusCode != res.StatusCode {
				t.Errorf("Unexpected statuc code, expected: %d, actual: %d", tt.wantStatusCode, res.StatusCode)
				return
			}

			if !strings.Contains(html, tt.wantTextInResHTML) {
				t.Errorf(`Response HTML does not contain "%s"`, tt.wantTextInResHTML)
				return
			}
		})
	}
}
