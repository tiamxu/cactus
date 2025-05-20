package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	defaultSigningKey  = "cactusOps"
	defaultExpiresIn   = time.Hour * 2
	defaultRefreshTime = time.Minute * 5
)

// 一些常量
var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("malformed token")
	ErrTokenInvalid     = errors.New("invalid token")
	ErrSigningKeyEmpty  = errors.New("signing key is empty")
)

type CustomClaims struct {
	UID int `json:"uid"`
	jwt.RegisteredClaims
}

type JWT struct {
	signingKey []byte
}

// NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	key := os.Getenv("JWT_SIGNING_KEY")
	if key == "" {
		key = defaultSigningKey
	}
	return &JWT{
		signingKey: []byte(key),
	}
}

// GenerateToken 生成令牌
func GenerateToken(uid int) (string, error) {
	claims := CustomClaims{
		UID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultExpiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(defaultSigningKey))
}

// ParseToken 解析并验证 Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(defaultSigningKey), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
		return nil, ErrTokenInvalid
	}
	if token == nil {
		return nil, ErrTokenInvalid
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, ErrTokenInvalid
}

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return GenerateToken(claims.UID)
}
