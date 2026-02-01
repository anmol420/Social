package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJwtAuthenticator(secret, aud, iss string) *JwtAuthenticator {
	return &JwtAuthenticator{
		secret: secret,
		aud:    aud,
		iss:    iss,
	}
}

func (a *JwtAuthenticator) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JwtAuthenticator) ValidateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
}
