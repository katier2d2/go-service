package auth

import (
	"crypto/rsa"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	RoleAdmin = "ADMIN"
	RoleUser  = "USER"
)

type ctxKey int

// Key is used to store/retrieve a Claims value from a context.Context
const Key ctxKey = 1

type Claims struct {
	jwt.StandardClaims
	Roles []string `json:"roles"`
}

func (c Claims) Authorize(roles ...string) bool {
	for _, has := range c.Roles {
		for _, want := range roles {
			if has == want {
				return true
			}
		}
	}
	return false
}

// in memory store of keys
type Keys map[string]*rsa.PrivateKey

// find the correct public key with a given kid
type PublicKeyLookup func(kid string) (*rsa.PublicKey, error)

// used to authenticate clients. can generate a token
type Auth struct {
	algorithm string
	keyFunc   func(t *jwt.Token) (interface{}, error)
	parser    *jwt.Parser
	keys      Keys
}

func New(algorithm string, lookup PublicKeyLookup, keys Keys) (*Auth, error) {
	if jwt.GetSigningMethod(algorithm) == nil {
		return nil, errors.Errorf("unknown algorithm %v", algorithm)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		kid, ok := t.Header["kid"]
		if !ok {
			return nil, errors.New("missing key id in token header")
		}
		kidId, ok := kid.(string)
		if !ok {
			return nil, errors.New("user token key id must be string")
		}
		return lookup(kidId)
	}

	parser := jwt.Parser{
		ValidMethods: []string{algorithm},
	}

	return &Auth{
		algorithm: algorithm,
		keyFunc:   keyFunc,
		parser:    &parser,
		keys:      keys,
	}, nil
}

func (a *Auth) AddKey(privateKey *rsa.PrivateKey, kid string) {
	a.keys[kid] = privateKey
}

func (a *Auth) RemoveKey(kid string) {
	delete(a.keys, kid)
}

func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	method := jwt.GetSigningMethod(a.algorithm)

	token := jwt.NewWithClaims(method, claims)
	token.Header["kid"] = kid

	privateKey, ok := a.keys[kid]
	if !ok {
		return "", errors.New("kid lookup failed")
	}

	str, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signing token")
	}

	return str, nil
}

func (a *Auth) ValidateToken(tokenStr string) (Claims, error) {
	var claims Claims
	token, err := a.parser.ParseWithClaims(tokenStr, &claims, a.keyFunc)
	if err != nil {
		return Claims{}, errors.Wrap(err, "parsing token")
	}

	if !token.Valid {
		return Claims{}, errors.Wrap(err, "invalid token")
	}

	return claims, nil
}
