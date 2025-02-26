// internal/middleware/auth.go

package middleware

import (
	"campus-wall/pkg/auth"
	"campus-wall/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
    AuthorizationHeaderKey  = "Authorization"
    AuthorizationTypeBearer = "Bearer"
    AuthorizationPayloadKey = "authorization_payload"
)

// @Summary 认证中间件
// @Description 验证用户token
// @Tags middleware
// @Param Authorization header string true "Bearer {token}"
func AuthMiddleware(tokenMaker *auth.JWTMaker) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader(AuthorizationHeaderKey)
        if len(authHeader) == 0 {
            err := errors.NewUnauthorizedError("authorization header is not provided")
            c.Error(err)
            c.Abort()
            return
        }

        fields := strings.Fields(authHeader)
        if len(fields) < 2 {
            err := errors.NewUnauthorizedError("invalid authorization header format")
            c.Error(err)
            c.Abort()
            return
        }

        authType := strings.ToLower(fields[0])
        if authType != strings.ToLower(AuthorizationTypeBearer) {
            err := errors.NewUnauthorizedError("unsupported authorization type")
            c.Error(err)
            c.Abort()
            return
        }

        accessToken := fields[1]
        claims, err := tokenMaker.VerifyToken(accessToken)
        if err != nil {
            err := errors.NewUnauthorizedError("invalid token")
            c.Error(err)
            c.Abort()
            return
        }

        c.Set(AuthorizationPayloadKey, claims)
        c.Next()
    }
}