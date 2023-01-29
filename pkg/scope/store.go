package scope

import (
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const symKey = "hoge"

type ScopeStore map[string]string

type ScopeStoreError struct {
	msg string
	err error
}

func (e *ScopeStoreError) Error() string {
	return fmt.Sprintf("cannot store error: %s (%s)", e.msg, e.err.Error())
}

func (e *ScopeStoreError) Unwrap() error {
	return e.err
}

func (ss ScopeStore) Get(tok string) (string, error) {
	if s, ok := ss[tok]; ok {
		return s, nil
	}
	key, err := jwk.FromRaw([]byte(symKey))
	if err != nil {
		return "", &ScopeStoreError{ msg: "create key", err: err }
	}
	verifiedTok, err := jwt.Parse([]byte(tok), jwt.WithKey(jwa.HS256, key))
	if err != nil {
		return "", &ScopeStoreError{ msg: "parse token", err: err }
	}
	scope, ok := verifiedTok.Get("scope")
	if !ok {
		return "", &ScopeStoreError{ msg: "invalid token", err: nil }
	}
	if strScope, ok := scope.(string); ok {
		ss[tok] = strScope
		return strScope, nil
	}
	return "", &ScopeStoreError{ msg: "invalid scope", err: nil }
}

var store ScopeStore
