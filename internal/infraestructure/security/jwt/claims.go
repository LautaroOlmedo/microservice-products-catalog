package jwt

import "time"

// Claims representa los datos que viajan en el JWT.
// No es un User, no es dominio, es metadata de request.
type Claims struct {
	Scope     string    `json:"scope"`
	RequestID string    `json:"request_id"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}
