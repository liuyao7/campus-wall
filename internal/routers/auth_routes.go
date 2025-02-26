// internal/router/auth_routes.go

package router

import (
	"campus-wall/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (r *Router) setupAuthRoutes(api *gin.RouterGroup) {
    auth := api.Group("/").Use(middleware.AuthMiddleware(r.cfg.TokenMaker))
    
    r.setupUserRoutes(auth)
    r.setupPostRoutes(auth)
    r.setupLikeRoutes(auth)
    r.setupCommentRoutes(auth)
}

func (r *Router) setupUserRoutes(auth *gin.RouterGroup) {
    if r.cfg.UserHandler == nil {
        return
    }

    users := auth.Group("/users")
    {
        users.PUT("/me", r.cfg.UserHandler.UpdateUser)
        users.POST("/avatar", r.cfg.UserHandler.UploadAvatar)
    }
}

func (r *Router) setupPostRoutes(auth *gin.RouterGroup) {
    if r.cfg.PostHandler == nil {
        return
    }

    posts := auth.Group("/posts")
    {
        posts.POST("/images", r.cfg.PostHandler.UploadImages)
        posts.POST("", r.cfg.PostHandler.CreatePost)
        posts.GET("", r.cfg.PostHandler.GetPostList)
        posts.GET("/users/:user_id", r.cfg.PostHandler.GetUserPosts)
        posts.GET("/topics/:topic_id", r.cfg.PostHandler.GetTopicPosts)
    }
}

// ... 类似地实现点赞和评论路由