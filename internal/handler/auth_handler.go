package handler

import (
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
    Code     string `json:"code"`
    Nickname string `json:"nickname"`
    Avatar   string `json:"avatar"`
    Gender   int    `json:"gender"`
}

func (h *Handler) WXLogin(c *gin.Context) {
    // var req LoginRequest
    // if err := c.ShouldBindJSON(&req); err != nil {
    //     c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
    //     return
    // }

    // // 调用微信登录接口
    // wxResp, err := service.WXLogin(req.Code)
    // if err != nil {
    //     c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
    //     return
    // }

    // // 查找或创建用户
    // user, err := h.repo.GetUserByOpenID(wxResp.OpenID)
    // if err != nil {
    //     // 用户不存在，创建新用户
    //     newUser := &model.User{
    //         OpenID:   wxResp.OpenID,
    //         Nickname: req.Nickname,
    //         Avatar:   req.Avatar,
    //         Gender:   req.Gender,
    //     }
    //     if err := h.repo.CreateUser(newUser); err != nil {
    //         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
    //         return
    //     }
    //     user = newUser
    // }

    // // 生成JWT token
    // token, err := auth.CreateToken(user.ID, user.OpenID)
    // if err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
    //     return
    // }

    // c.JSON(http.StatusOK, gin.H{
    //     "token": token,
    //     "user":  user,
    // })
}