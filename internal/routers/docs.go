// internal/router/docs.go

package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (r *Router) setupDocs(api *gin.RouterGroup) {
    // 仅在非生产环境启用 swagger
    if gin.Mode() != gin.ReleaseMode {
        api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
        
        // API 文档主页重定向到 swagger
        api.GET("/docs", func(c *gin.Context) {
            c.Redirect(301, "/api/swagger/index.html")
        })
    }
}