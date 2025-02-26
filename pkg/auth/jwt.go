package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
    ErrInvalidToken = errors.New("invalid token")
    ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
    UserID uint64 `json:"user_id"`
    OpenID string `json:"open_id"`
    jwt.RegisteredClaims
}

type JWTMaker struct {
    secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
    return &JWTMaker{secretKey: secretKey}
}

// CreateToken 创建一个新的token
func (maker *JWTMaker) CreateToken(userID uint64, openID string) (string, error) {
    claims := &Claims{
        UserID: userID,
        OpenID: openID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(maker.secretKey))
}

// VerifyToken 验证token
func (maker *JWTMaker) VerifyToken(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenStr,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(maker.secretKey), nil
        },
    )

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return nil, ErrInvalidToken
        }
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, ErrInvalidToken
    }

    if time.Now().After(claims.ExpiresAt.Time) {
        return nil, ErrExpiredToken
    }

    return claims, nil
}