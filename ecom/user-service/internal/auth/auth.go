package auth

import (
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey int

const ClaimsKey ctxKey = 1

type Keys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}
type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// NewKeys is a constructor function for Keys struct. It accepts privateKey and publicKey as parameters and returns
// an instance of Keys struct. If either of privateKey or publicKey is nil, it returns an error.
func NewKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (*Keys, error) {
	if privateKey == nil || publicKey == nil {
		return nil, fmt.Errorf("invalid keys")
	}
	return &Keys{privateKey, publicKey}, nil
}

// GenerateToken is a method for Auth struct. It generates a new JWT token using the provided claims and
// signs it using the privateKey of the Auth struct it's called upon. If there is an error during signing,
// it returns an error.
func (k *Keys) GenerateToken(claims Claims) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenStr, err := tkn.SignedString(k.privateKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ValidateToken is a method for Auth struct. It verifies the provided JWT token using the publicKey of the Auth struct
// it's called upon and returns the parsed claims if the JWT token is valid. If the JWT token is invalid or
// there is an error during parsing, it returns an error.
func (k *Keys) ValidateToken(tokenStr string) (Claims, error) {
	var claims Claims
	tkn, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return k.publicKey, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !tkn.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}
	return claims, nil
}
