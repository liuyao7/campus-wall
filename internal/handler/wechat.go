// internal/handler/wechat.go

package handler

import (
	"campus-wall/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WeChatHandler struct {
    wechatService *service.WeChatService
}

func NewWeChatHandler(wechatService *service.WeChatService) *WeChatHandler {
    return &WeChatHandler{
        wechatService: wechatService,
    }
}

// @Summary 微信小程序登录
// @Description 通过微信小程序登录获取token
// @Tags 微信
// @Accept json
// @Produce json
// @Param code body string true "微信登录code"
// @Success 200 {object} response.Response{data=auth.TokenResponse} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /wx/login [post]
func (h *WeChatHandler) MiniProgramLogin(c *gin.Context) {
    var req service.WXLoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    fmt.Println(req)

    response, err := h.wechatService.MiniProgramLogin(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// @Summary 解密用户信息
// @Description 解密微信用户信息
// @Tags 微信
// @Accept json
// @Produce json
// @Param request body wechat.DecryptUserInfoRequest true "解密请求参数"
// @Success 200 {object} response.Response{data=wechat.UserInfo} "解密成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /wx/decrypt/userinfo [post]
func (h *WeChatHandler) DecryptUserInfo(c *gin.Context) {
    var req service.WXDecryptRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userInfo, err := h.wechatService.DecryptUserInfo(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, userInfo)
}

func (h *WeChatHandler) DecryptPhoneNumber(c *gin.Context) {
    var req service.WXDecryptRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    phoneInfo, err := h.wechatService.DecryptPhoneNumber(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, phoneInfo)
}