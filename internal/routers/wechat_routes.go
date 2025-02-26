// internal/router/wechat_routes.go

package router

import "github.com/gin-gonic/gin"

func (r *Router) setupWeChatRoutes(api *gin.RouterGroup) {
    if r.cfg.WechatHandler == nil {
        return
    }

    wx := api.Group("/wx")
    {
        wx.POST("/login", r.cfg.WechatHandler.MiniProgramLogin)
        wx.POST("/decrypt/userinfo", r.cfg.WechatHandler.DecryptUserInfo)
        wx.POST("/decrypt/phone", r.cfg.WechatHandler.DecryptPhoneNumber)
    }
}