// internal/handler/user.go

package handler

import (
	"campus-wall/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
    userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
    }
}

func (h *UserHandler) Register(c *gin.Context) {
    var req service.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.userService.Register(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Login(c *gin.Context) {
    var req service.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.userService.Login(&req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// UpdateUser 更新用户信息
// func (h *UserHandler) UpdateUser(c *gin.Context) {
//     // 1. 获取当前用户ID
//     userID, exists := c.Get(middleware.UserIDKey)
//     if !exists {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
//         return
//     }

//     // 2. 验证请求数据
//     var req model.UserUpdateRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // 3. 更新用户信息
//     updatedUser, err := h.userService.UpdateUser(userID.(uint), &req)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     // 4. 返回更新后的用户信息
//     c.JSON(http.StatusOK, updatedUser)
// }

// func (h *UserHandler) UploadAvatar(c *gin.Context) {
//     // 1. 获取当前用户ID
//     userID, exists := c.Get(middleware.UserIDKey)
//     if !exists {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
//         return
//     }

//     // 2. 获取上传的文件
//     file, err := c.FormFile("avatar")
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
//         return
//     }

//     // 3. 上传头像
//     user, err := h.userService.UploadAvatar(userID.(uint), file)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, user)
// }