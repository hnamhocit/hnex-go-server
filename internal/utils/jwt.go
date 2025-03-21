package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Sub uint `json:"sub"`
	jwt.RegisteredClaims
}

type Config struct {
	AccessSecret     string
	RefreshSecret    string
	AccessExpiresIn  time.Time
	RefreshExpiresIn time.Time
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GetJWTConfig() *Config {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	accessExpiresIn := time.Now().Add(1 * time.Hour)
	refreshExpiresIn := time.Now().Add(7 * 24 * time.Hour)

	config := &Config{
		AccessSecret:     accessSecret,
		RefreshSecret:    refreshSecret,
		AccessExpiresIn:  accessExpiresIn,
		RefreshExpiresIn: refreshExpiresIn,
	}

	return config
}

func GenerateTokens(sub uint) (*Tokens, error) {
	config := GetJWTConfig()

	accessClaims := &Claims{
		Sub: sub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(config.AccessExpiresIn),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshClaims := &Claims{
		Sub: sub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(config.RefreshExpiresIn),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessTokenString, err := accessToken.SignedString([]byte(config.AccessSecret))
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString([]byte(config.RefreshSecret))
	if err != nil {
		return nil, err
	}

	tokens := &Tokens{AccessToken: accessTokenString, RefreshToken: refreshTokenString}
	return tokens, nil
}

func ValidateToken(tokenString, jwtKey string) (*Claims, error) {
	config := GetJWTConfig()
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if jwtKey == "JWT_ACCESS_SECRET" {
			return []byte(config.AccessSecret), nil
		}

		return []byte(config.RefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
