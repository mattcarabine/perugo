package server

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestJwtValidation(t *testing.T) {
	validToken, err := generateJwt()
	if err != nil {
		t.Error("automatic token generation failed")
	}
	var tests = []struct {
		descr string
		token string
		valid bool
	}{
		{"Happy Path", validToken, true},
		{"Expiry earlier than issued", func() string {
			now := time.Now()
			claims := jwt.StandardClaims{
				ExpiresAt: 15000,
				IssuedAt:  now.Unix(),
				NotBefore: now.Unix(),
				Issuer:    "perugo-api",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// Sign and get the complete encoded token as a string using the secret
			tokenStr, _ := token.SignedString(secret)
			return tokenStr
		}(), false},
		{"Expiry earlier than now", func() string {
			now := time.Now()
			claims := jwt.StandardClaims{
				ExpiresAt: 15000,
				NotBefore: now.Unix(),
				Issuer:    "perugo-api",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			// Sign and get the complete encoded token as a string using the secret
			tokenStr, _ := token.SignedString(secret)
			return tokenStr
		}(), false},
		{"Token not yet valid", func() string {
			now := time.Now()
			claims := jwt.StandardClaims{
				ExpiresAt: now.Add(5 * time.Minute).Unix(),
				IssuedAt:  now.Unix(),
				NotBefore: now.Add(3 * time.Minute).Unix(),
				Issuer:    "perugo-api",
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenStr, _ := token.SignedString(secret)
			return tokenStr
		}(), false},
		{"Junk token", "testing", false},
	}

	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			valid, err := validateJWT(tt.token)
			if valid != tt.valid {
				t.Errorf("expected valid %v, got %v, err: %v", tt.valid, valid, err)
			}
		})
	}

}
