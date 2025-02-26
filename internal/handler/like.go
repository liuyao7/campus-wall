// // internal/handler/like.go

package handler

// import (
// 	"campus-wall/internal/middleware"
// 	"campus-wall/internal/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type LikeHandler struct {
//     likeService *service.LikeService
// }

// func NewLikeHandler(likeService *service.LikeService) *LikeHandler {
//     return &LikeHandler{
//         likeService: likeService,
//     }
// }

// // LikePost 点赞帖子
// func (h *LikeHandler) LikePost(c *gin.Context) {
//     userID, _ := c.Get(middleware.UserIDKey)
//     postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
//         return
//     }

//     if err := h.likeService.LikePost(userID.(uint), uint(postID)); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "success"})
// }

// // UnlikePost 取消点赞
// func (h *LikeHandler) UnlikePost(c *gin.Context) {
//     userID, _ := c.Get(middleware.UserIDKey)
//     postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
//         return
//     }

//     if err := h.likeService.UnlikePost(userID.(uint), uint(postID)); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "success"})
// }

// // GetPostLikeUsers 获取帖子点赞用户列表
// func (h *LikeHandler) GetPostLikeUsers(c *gin.Context) {
//     postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
//         return
//     }

//     page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//     size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

//     users, total, err := h.likeService.GetPostLikeUsers(uint(postID), page, size)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{
//         "users": users,
//         "total": total,
//         "page":  page,
//         "size":  size,
//     })
// }

// // GetUserLikedPosts 获取用户点赞的帖子列表
// func (h *LikeHandler) GetUserLikedPosts(c *gin.Context) {
//     userID, _ := c.Get(middleware.UserIDKey)
//     page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//     size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

//     posts, total, err := h.likeService.GetUserLikedPosts(userID.(uint), page, size)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{
//         "posts": posts,
//         "total": total,
//         "page":  page,
//         "size":  size,
//     })
// }