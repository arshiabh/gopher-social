package auth

import (
	"fmt"
	"log"

	"github.com/golang-jwt/jwt/v5"
)

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
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(a.secret), nil
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["sub"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
	return jwtToken, err
}
