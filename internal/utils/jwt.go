package utils

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenClaims struct {
	Sub uint `json:"sub"`
	jwt.RegisteredClaims
}

func GenerateTokens(sub uint) (Tokens, error) {
	accessTokenSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshTokenSecret := os.Getenv("JWT_REFRESH_SECRET")

	if accessTokenSecret == "" || refreshTokenSecret == "" {
		return Tokens{}, errors.New("missing token secret in environment variables")
	}

	accessTokenClaims := TokenClaims{
		Sub: sub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	refreshTokenClaims := TokenClaims{
		Sub: sub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 7 * 24)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessTokenString, err := accessToken.SignedString([]byte(accessTokenSecret))
	if err != nil {
		return Tokens{}, fmt.Errorf("error generating access token: %w", err)
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return Tokens{}, fmt.Errorf("error generating refresh token: %w", err)
	}

	return Tokens{AccessToken: accessTokenString, RefreshToken: refreshTokenString}, nil
}

func VerifyToken(s, secretKey string) (*jwt.Token, error) {
	secret := os.Getenv(secretKey)
	if secret == "" {
		return nil, fmt.Errorf("missing secret key for token verification")
	}

	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token, nil
}

func IsTokenExpired(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(float64)

	if !ok {
		return true
	}

	expirationTime := time.Unix(int64(exp), 0)

	return time.Now().After(expirationTime)
}
