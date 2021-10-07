package auth_test

import (
	"calendar/internal/app/httpserver/auth"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	secretKey = "secretkey"
	issuer    = "AuthService"
	login     = "test"
)

func TestJwtWrapper_GenerateToken(t *testing.T) {
	jwtWrapper := auth.JwtWrapper{
		SecretKey:       secretKey,
		Issuer:          issuer,
		ExpirationHours: 24,
	}

	generateToken, err := jwtWrapper.GenerateToken(login)
	assert.NoError(t, err)

	os.Setenv("testToken", generateToken)
}

func TestJwtWrapper_ValidateToken(t *testing.T) {
	encodedToken := os.Getenv("testToken")

	jwtWrapper := auth.JwtWrapper{
		SecretKey: secretKey,
		Issuer:    issuer,
	}

	claims, err := jwtWrapper.ValidateToken(encodedToken)
	assert.NoError(t, err)

	assert.Equal(t, login, claims.Login)
	assert.Equal(t, issuer, claims.Issuer)

}
