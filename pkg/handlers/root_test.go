package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/osemisan/osemisan-client/pkg/handlers"
	"github.com/osemisan/osemisan-client/testutil"
)

func TestRootHandler(t *testing.T) {
	c := new(http.Client)

	tests := []struct {
		name           string
		token          string
		withoutTok     bool
		wantStatusCode int
	}{
		{
			"トークンがクッキーに格納されていないとき、ステータスコード200",
			"",
			true,
			http.StatusOK,
		},
		{
			"正しいトークンがクッキーに格納されているとき、ステータスコード200",
			testutil.BuildScopedJwt(t, testutil.Scopes{}),
			false,
			http.StatusOK,
		},
		{
			"間違ったトークンがクッキーに格納されているとき、ステータスコード500",
			"invalidtoken",
			false,
			http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := httptest.NewServer(http.HandlerFunc(handlers.RootHandler))
			defer s.Close()

			req, err := http.NewRequest(http.MethodGet, s.URL, nil)
			if err != nil {
				t.Error("Failed to create new request", err)
				return
			}

			if !tt.withoutTok {
				req.AddCookie(&http.Cookie{
					Name:  "token",
					Value: tt.token,
				})
			}
			res, err := c.Do(req)
			if err != nil {
				t.Error("Failed to request", err)
				return
			}

			if res.StatusCode != tt.wantStatusCode {
				t.Errorf("Unexpected status code, expected: %d, actual: %d", tt.wantStatusCode, res.StatusCode)
				return
			}
		})
	}
}
