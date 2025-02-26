// internal/middleware/sensitive.go

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/importcjj/sensitive"
)

func SensitiveFilter() gin.HandlerFunc {
    filter := sensitive.New()
    filter.LoadWordDict("configs/sensitive_words.txt")

    return func(c *gin.Context) {
        var data map[string]interface{}
        if err := c.ShouldBindJSON(&data); err != nil {
            c.Next()
            return
        }

        if content, ok := data["content"].(string); ok {
            if found, _ := filter.Validate(content); found {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": "content contains sensitive words",
                })
                c.Abort()
                return
            }
            // 替换敏感词
            data["content"] = filter.Replace(content, '*')
        }

        c.Set("filtered_data", data)
        c.Next()
    }
}