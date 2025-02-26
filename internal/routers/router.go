// internal/router/router.go

package router

import (
	"campus-wall/internal/handler"
	"campus-wall/internal/middleware"
	"campus-wall/pkg/auth"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

type Router struct {
    engine *gin.Engine
    cfg    *Config
}

type Config struct {
    TokenMaker     *auth.JWTMaker
    WechatHandler  *handler.WeChatHandler
    UserHandler    *handler.UserHandler
    // PostHandler    *handler.PostHandler
    // LikeHandler    *handler.LikeHandler
    // CommentHandler *handler.CommentHandler
    RateLimit      rate.Limit
    RateBurst      int
}

// NewRouter 创建新的路由实例
func NewRouter(cfg *Config) *Router {
    r := &Router{
        engine: gin.New(),
        cfg:    cfg,
    }
    r.setupMiddlewares()
    r.setupRoutes()
    return r
}

// Engine 返回gin引擎实例
func (r *Router) Engine() *gin.Engine {
    return r.engine
}

// setupMiddlewares 设置中间件
func (r *Router) setupMiddlewares() {
    r.engine.Use(
        gin.Recovery(),
        middleware.Logger(),
        middleware.CORS(),
        middleware.RateLimit(r.cfg.RateLimit, r.cfg.RateBurst),
    )
}

// setupRoutes 设置所有路由
func (r *Router) setupRoutes() {
	// swagger
    r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    api := r.engine.Group("/api")
    
    r.setupPublicRoutes(api)
    r.setupWeChatRoutes(api)
    r.setupAuthRoutes(api)
}