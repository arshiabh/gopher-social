package auth

import "github.com/golang-jwt/jwt/v5"

type JWTAuth struct {
	secret string
}

func NewAuthentication(secret string) *JWTAuth {
	return &JWTAuth{
		secret: secret,
	}
}

func (a *JWTAuth) GenerateToken(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *JWTAuth) ValidateToken(token string) (*jwt.Token, error) {
	return nil, nil
}
