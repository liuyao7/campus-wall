// internal/middleware/logger.go

package middleware

import (
	"campus-wall/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 开始时间
        start := time.Now()
        path := c.Request.URL.Path
        query := c.Request.URL.RawQuery

        // 处理请求
        c.Next()

        // 结束时间
        end := time.Now()
        latency := end.Sub(start)

        // 记录日志
        logger.Info("HTTP Request",
            zap.Int("status", c.Writer.Status()),
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", query),
            zap.String("ip", c.ClientIP()),
            zap.String("user-agent", c.Request.UserAgent()),
            zap.Duration("latency", latency),
            zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
        )
    }
}