package scope_test

import (
	"testing"

	"github.com/osemisan/osemisan-client/pkg/scope"
	"github.com/osemisan/osemisan-client/testutil"
)

func TestGet(t *testing.T) {
	tests := []struct{
		name string
		repeat int
		token string
		wantScope string
	}{
		{
			"格納したトークンに対応するスコープが返ってくる(1回目)",
			1,
			testutil.BuildScopedJwt(t, testutil.Scopes{
				Abura: true,
				Minmin: true,
				Kuma: false,
				Niinii: false,
				Tsukutsuku: false,
			}),
			"abura minmin",
		},
		{
			"格納したトークンに対応するスコープが返ってくる(2回目)",
			2,
			testutil.BuildScopedJwt(t, testutil.Scopes{
				Abura: true,
				Minmin: true,
				Kuma: true,
				Niinii: false,
				Tsukutsuku: false,
			}),
			"abura minmin kuma",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := scope.ScopeStore{}
			var got string
			for i := 0; i < tt.repeat; i++ {
				scope, err := store.Get(tt.token)
				if err != nil {
					t.Error("Failed to get token from store", err)
					return
				}
				got = scope
			}
			if got != tt.wantScope {
				t.Errorf(`Unexpected scope, expected: "%s", actual: "%s"`, tt.wantScope, got)
			}
		})
	}
}
