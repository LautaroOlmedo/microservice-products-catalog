package jwt

import (
	"fmt"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"time"
)

type Verifier interface {
	Verify(token string) (Claims, error)
}

type JWTVerifier struct {
	secret []byte
}

func NewVerifier(secret string) *JWTVerifier {
	return &JWTVerifier{
		secret: []byte(secret),
	}
}

func (v *JWTVerifier) Verify(tokenString string) (Claims, error) {
	token, err := jwtlib.Parse(tokenString, func(token *jwtlib.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return v.secret, nil
	})

	if err != nil || !token.Valid {
		return Claims{}, fmt.Errorf("invalid token")
	}

	mapClaims, ok := token.Claims.(jwtlib.MapClaims)
	if !ok {
		return Claims{}, fmt.Errorf("invalid claims")
	}

	return Claims{
		Scope:     mapClaims["scope"].(string),
		RequestID: mapClaims["request_id"].(string),
		IssuedAt:  time.Unix(int64(mapClaims["iat"].(float64)), 0),
		ExpiresAt: time.Unix(int64(mapClaims["exp"].(float64)), 0),
	}, nil
}
