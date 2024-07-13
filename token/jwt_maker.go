package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKiySize = 32

// * Note [codermuss]: JWTMaker is a JSON Web Token Maker
type JWTMaker struct {
	secretKey string
}

// * Note [codermuss]: NewJWTMaker creates a new JWTMaker

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKiySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKiySize)
	}
	return &JWTMaker{secretKey}, nil
}

// * Note [codermuss]: CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, role string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, role, duration)
	if err != nil {
		return "", payload, err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// * Note [codermuss]: VerifyToken checks if the token is valid or not

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretKey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
