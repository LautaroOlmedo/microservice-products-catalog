package jwt

import "time"

type Claims struct {
	Scope     string    `json:"scope"`
	RequestID string    `json:"request_id"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}
