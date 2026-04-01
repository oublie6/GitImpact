// Package auth 提供 JWT 生成与解析能力。
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 是写入 JWT 的业务字段。
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 按给定密钥和过期时间生成 JWT 字符串。
func GenerateToken(secret string, ttl time.Duration, userID uint, username, role string) (string, error) {
	now := time.Now()
	claims := Claims{UserID: userID, Username: username, Role: role, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(now.Add(ttl)), IssuedAt: jwt.NewNumericDate(now)}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析并校验 JWT，返回业务 claims。
func ParseToken(secret, t string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(t, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
