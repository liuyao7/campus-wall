package handler

import (
	_ "campus-wall/internal/model"
	_ "net/http"

	"campus-wall/internal/repository"

	_ "github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
)

type Handler struct {
    // db *gorm.DB
	repo *repository.Repository
}

func NewHandler(repo *repository.Repository) *Handler {
    // return &Handler{db: db}
    return &Handler{repo: repo,}
}

// 创建帖子
// func (h *Handler) CreatePost(c *gin.Context) {
//     var post model.Post
//     if err := c.ShouldBindJSON(&post); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // TODO: 从JWT获取用户ID
//     // post.UserID = getUserIDFromToken(c)

//     if err := h.db.Create(&post).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, post)
// }

// 获取帖子列表
// func (h *Handler) GetPosts(c *gin.Context) {
//     var posts []model.Post
//     if err := h.db.Preload("User").Find(&posts).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, posts)
// }

// 获取帖子详情
// func (h *Handler) GetPost(c *gin.Context) {
// 	var post model.Post
// 	if err := h.db.Preload("User").First(&post, c.Param("id")).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, post)
// }

// 添加评论
// func (h *Handler) CreateComment(c *gin.Context) {
//     var comment model.Comment
//     if err := c.ShouldBindJSON(&comment); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // TODO: 从JWT获取用户ID
//     // comment.UserID = getUserIDFromToken(c)

//     if err := h.db.Create(&comment).Error; err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, comment)
// }

// 获取评论列表
// func (h *Handler) GetComments(c *gin.Context) {
// 	var comments []model.Comment
// 	if err := h.db.Preload("User").Find(&comments).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, comments)
// }

// CreateLike 创建点赞
// func (h *Handler) CreateLike(c *gin.Context) {
// 	var like model.Like
// 	if err := c.ShouldBindJSON(&like); err != nil {
// 					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 					return
// 	}

// 	// TODO: 从JWT获取用户ID
// 	// like.UserID = getUserIDFromToken(c)

// 	if err := h.db.Create(&like).Error; err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 					return
// 	}

// 	c.JSON(http.StatusOK, like)
// }

// DeleteLike 删除点赞
// func (h *Handler) DeleteLike(c *gin.Context) {
// 	var like model.Like
// 	if err := c.ShouldBindJSON(&like); err != nil {
// 					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 					return
// 	}

// 	// TODO: 从JWT获取用户ID
// 	// like.UserID = getUserIDFromToken(c)

// 	if err := h.db.Delete(&like).Error; err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 					return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "点赞已删除"})
// }
