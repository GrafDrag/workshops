package auth_test

import (
	"calendar/internal/app/httpserver/auth"
	"calendar/internal/model"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	secretKey = "secretkey"
	issuer    = "AuthService"
	userID    = 1
)

func TestJwtWrapper_GenerateToken(t *testing.T) {
	jwtWrapper := auth.JwtWrapper{
		SecretKey:       secretKey,
		Issuer:          issuer,
		ExpirationHours: 24,
	}
	u := model.TestUser(t)
	u.ID = userID

	generateToken, err := jwtWrapper.GenerateToken(u)
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

	assert.Equal(t, userID, claims.ID)
	assert.Equal(t, issuer, claims.Issuer)
}
