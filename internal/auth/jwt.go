package auth

import "github.com/golang-jwt/jwt/v5"

type JwtAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJwtAuthenticator(secret, aud, iss string) *JwtAuthenticator {
	return &JwtAuthenticator{
		secret: secret,
		aud: aud,
		iss: iss,
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
	return nil, nil
}
