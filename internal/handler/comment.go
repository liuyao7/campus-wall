// // internal/handler/comment.go

package handler

// import (
// 	"campus-wall/internal/model"
// 	"campus-wall/internal/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type CommentHandler struct {
//     commentService *service.CommentService
// }

// func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
//     return &CommentHandler{
//         commentService: commentService,
//     }
// }

// // CreateComment 创建评论
// func (h *CommentHandler) CreateComment(c *gin.Context) {
//     userID := c.MustGet("user_id").(uint)
//     postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
//         return
//     }

//     var req model.CommentRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     comment, err := h.commentService.CreateComment(userID, &req, uint(postID))
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, comment)
// }

// // DeleteComment 删除评论
// func (h *CommentHandler) DeleteComment(c *gin.Context) {
//     userID := c.MustGet("user_id").(uint)
//     commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
//         return
//     }

//     if err := h.commentService.DeleteComment(userID, uint(commentID)); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "success"})
// }

// // GetCommentList 获取评论列表
// func (h *CommentHandler) GetCommentList(c *gin.Context) {
//     var query model.CommentListQuery
//     if err := c.ShouldBindQuery(&query); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     var userID *uint
//     if id, exists := c.Get("user_id"); exists {
//         uid := id.(uint)
//         userID = &uid
//     }

//     response, err := h.commentService.GetCommentList(&query, userID)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, response)
// }

// // LikeComment 点赞评论
// func (h *CommentHandler) LikeComment(c *gin.Context) {
//     userID := c.MustGet("user_id").(uint)
//     commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
//         return
//     }

//     if err := h.commentService.LikeComment(userID, uint(commentID)); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "success"})
// }

// // UnlikeComment 取消点赞评论
// func (h *CommentHandler) UnlikeComment(c *gin.Context) {
//     userID := c.MustGet("user_id").(uint)
//     commentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
//         return
//     }

//     if err := h.commentService.UnlikeComment(userID, uint(commentID)); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "success"})
// }