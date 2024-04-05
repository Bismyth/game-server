package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Id uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

var signingKey []byte = []byte("secretKey")

// Override signing key from config
func SetSigningKey(s string) {
	signingKey = []byte(s)
}

func VerifyToken(t string) (uuid.UUID, error) {
	var claims Claims

	if t == "" {
		return uuid.Nil, nil
	}

	token, err := jwt.ParseWithClaims(string(t), &claims, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return uuid.Nil, nil
		}
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, nil
	}

	return claims.Id, nil
}

const expireTime time.Duration = time.Hour * 3

func GenerateToken(userId uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(expireTime)

	claims := Claims{
		Id: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
