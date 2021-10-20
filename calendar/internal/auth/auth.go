package auth

import (
	"calendar/internal/config"
	"calendar/internal/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type ctxUserKey int

const (
	KeyUserID ctxUserKey = iota
)

type JwtWrapper struct {
	SecretKey       string
	ExpirationHours int64
	Issuer          string
}

func NewJwtWrapper(config config.JWTConfig, issuer string) *JwtWrapper {
	return &JwtWrapper{
		SecretKey:       config.JwtSecretKey,
		ExpirationHours: config.JwtExpHours,
		Issuer:          issuer,
	}
}

type JwtClaim struct {
	ID int
	jwt.StandardClaims
}

func (j *JwtWrapper) GenerateToken(user *model.User) (string, error) {
	claims := &JwtClaim{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(j.ExpirationHours)).Unix(),
			Issuer:    j.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *JwtWrapper) ValidateToken(signedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.SecretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		return nil, errors.New("couldn't parse claims")
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}
