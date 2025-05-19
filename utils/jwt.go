package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	defaultSigningKey  = "cactusOps"     // 生产环境必须修改！
	defaultExpiresIn   = time.Hour * 24  // 默认 24 小时
	defaultRefreshTime = time.Minute * 5 // 剩余 5 分钟时刷新
)

// 一些常量
var (
	ErrTokenExpired     = errors.New("token expired")
	ErrTokenNotValidYet = errors.New("token not active yet")
	ErrTokenMalformed   = errors.New("malformed token")
	ErrTokenInvalid     = errors.New("invalid token")
	ErrSigningKeyEmpty  = errors.New("signing key is empty")
)

// CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	UID int `json:"uid"`
	jwt.RegisteredClaims
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte `json:"signing_key"`
}

// NewJWT 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(os.Getenv("JWT_SIGNING_KEY")),
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
