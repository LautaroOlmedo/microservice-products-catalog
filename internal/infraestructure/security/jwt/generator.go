// internal/infrastructure/security/jwt/generator.go
package jwt

import (
	"context"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"microservice-products-catalog/cmd/http/auth"
)

type JWTGenerator struct {
	secret []byte
	ttl    time.Duration
}

func (g *JWTGenerator) Generate(
	ctx context.Context,
	input auth.TokenClaims,
) (string, error) {

	now := time.Now()

	claims := jwtlib.MapClaims{
		"scope":      input.Scope,
		"request_id": input.RequestID,
		"iat":        now.Unix(),
		"exp":        now.Add(g.ttl).Unix(),
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(g.secret)
}

func NewTokenGenerator(secret string, ttl time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secret: []byte(secret),
		ttl:    ttl,
	}
}
