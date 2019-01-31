package jwt

import (
	"bytes"
	"encoding/json"

	"cirello.io/errors"
	jwt "github.com/dgrijalva/jwt-go"
)

// Parse decodes the JWT from the given string. It will return only a valid
// token, and an error otherwise.
func Parse(t string, caPEM []byte) (*jwt.Token, ServiceClaims, error) {
	var claims ServiceClaims
	token, err := jwt.ParseWithClaims(t, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return caPEM, nil
		})
	if err != nil {
		return nil, ServiceClaims{},
			errors.E(errors.Invalid, err, "cannot parse token")
	}
	if !token.Valid {
		return nil, ServiceClaims{},
			errors.E(errors.Invalid, "not is not valid")
	}
	return token, claims, nil
}

// Claims from a given token. It will return not OK if a ServiceClaim is not
// found.
func Claims(t *jwt.Token) (ServiceClaims, error) {
	var sc ServiceClaims
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return sc, errors.E("not a claim map")
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(claims); err != nil {
		return sc, errors.E(err, "cannot encode claim map")
	}
	if err := json.NewDecoder(&buf).Decode(&sc); err != nil {
		return sc, errors.E(err, "cannot decode service claims")
	}
	return sc, nil
}
