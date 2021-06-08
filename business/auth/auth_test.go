package auth_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/katier2d2/go-service/business/auth"
	"github.com/stretchr/testify/assert"
)

const (
	success = "\u2713"
	failure = "\u2717"
)

func TestAuth(t *testing.T) {
	t.Log("Given the need to be able to authenticate and authorize access.")
	{
		testId := 0
		t.Logf("\tTest %d:\tWhen handling a single user.", testId)
		{
			privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
			assert.NoError(t, err)

			const keyId = "64651844"
			lookup := func(kid string) (*rsa.PublicKey, error) {
				switch kid {
				case keyId:
					return &privateKey.PublicKey, nil
				}
				return nil, fmt.Errorf("no public key found for the specified kid: %s", kid)
			}

			a, err := auth.New("RS256", lookup, auth.Keys{keyId: privateKey})
			assert.NoError(t, err)

			claims := auth.Claims{
				StandardClaims: jwt.StandardClaims{
					Issuer:    "service project",
					Subject:   "64651844",
					Audience:  "students",
					ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
					IssuedAt:  time.Now().Unix(),
				},
				Roles: []string{auth.RoleAdmin},
			}

			token, err := a.GenerateToken(keyId, claims)
			assert.NoError(t, err)

			parsedClaims, err := a.ValidateToken(token)
			assert.NoError(t, err)

			if exp, got := len(claims.Roles), len(parsedClaims.Roles); exp != got {
				t.Logf("\tTest %d:\texp: %v", testId, exp)
				t.Logf("\tTest %d:\tgot: %v", testId, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected roles", failure, testId)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected roles", success, testId)

			if exp, got := len(claims.Roles[0]), len(parsedClaims.Roles[0]); exp != got {
				t.Logf("\tTest %d:\texp: %v", testId, exp)
				t.Logf("\tTest %d:\tgot: %v", testId, got)
				t.Fatalf("\t%s\tTest %d:\tShould have the expected roles", failure, testId)
			}
			t.Logf("\t%s\tTest %d:\tShould have the expected roles", success, testId)

		}
	}
}
