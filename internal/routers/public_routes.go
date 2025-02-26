// internal/router/public_routes.go

package router

import (
	"campus-wall/internal/handler"

	"github.com/gin-gonic/gin"
)


func (r *Router) setupPublicRoutes(api *gin.RouterGroup) {
    public := api.Group("/")
    {
        public.GET("/ping", handler.Ping)
        if r.cfg.UserHandler != nil {
            public.POST("/register", r.cfg.UserHandler.Register)
            public.POST("/login", r.cfg.UserHandler.Login)
        }
    }
}