package api

import (
	"fmt"
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

const expireTime time.Duration = time.Hour * 3

type RoomTokenClaims struct {
	RoomId uuid.UUID `json:"roomId"`
	UserId uuid.UUID `json:"userId"`
	jwt.RegisteredClaims
}

func GenerateRoomToken(roomId uuid.UUID, userId uuid.UUID) (string, error) {
	expirationTime := time.Now().Add(expireTime)

	claims := RoomTokenClaims{
		RoomId: roomId,
		UserId: userId,
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

func VerifyRoomToken(t string) (RoomTokenClaims, error) {
	var claims RoomTokenClaims

	if t == "" {
		return claims, fmt.Errorf("no room token given")
	}

	token, err := jwt.ParseWithClaims(string(t), &claims, func(t *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		return claims, fmt.Errorf("invalid room token")
	}

	return claims, nil
}
