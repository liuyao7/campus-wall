package handler

// import (
// 	"campus-wall/internal/model"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type CreatePostRequest struct {
//     Content  string `json:"content" binding:"required"`
//     Images   string `json:"images"`
//     Type     int    `json:"type"`
//     Location string `json:"location"`
// }

// func (h *Handler) CreatePost(c *gin.Context) {
//     var req CreatePostRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
//         return
//     }

//     userID := c.GetUint("userID")

//     post := &model.Post{
//         UserID:    userID,
//         Content:   req.Content,
//         Images:    req.Images,
//         Type:      req.Type,
//         Location:  req.Location,
//     }

//     if err := h.repo.CreatePost(post); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"data": post})
// }

// func (h *Handler) GetPosts(c *gin.Context) {
//     page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
//     pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

//     posts, total, err := h.repo.GetPosts(page, pageSize)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{
//         "data": posts,
//         "meta": gin.H{
//             "total": total,
//             "page": page,
//             "page_size": pageSize,
//         },
//     })
// }

// // 获取帖子详情
// func (h *Handler) GetPost(c *gin.Context) {
// 	postID := c.Param("id")
// 	if postID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 		return
// 	}

// 	postIDUint, err := strconv.Atoi(postID)

// 	post, err := h.repo.GetPostByID(uint(postIDUint))
// 	if err != nil {
// 					c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
// 					return
// 	}

// 	// 更新帖子浏览量
// 	err = h.repo.IncreaseViewCount(uint(postIDUint))
// 	if err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update views"})
// 					return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"data": post})
// }

// // 添加评论
// func (h *Handler) AddComment(c *gin.Context) {
// 	postID := c.Param("id")
// 	if postID == "" {
// 					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 					return
// 	}

// 	userID := c.GetUint("userID")
// 	content := c.PostForm("content")
// 	if content == "" {
// 					c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
// 					return
// 	}
// 	postIDUint, err := strconv.Atoi(postID)

// 	comment := &model.Comment{
// 					UserID:   userID,
// 					PostID:   uint(postIDUint),
// 					Content:  content,
// 	}

// 	err = h.repo.AddComment(comment)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"data": comment})
// }

// // 获取评论列表
// func (h *Handler) GetComments(c *gin.Context) {
// 	postID := c.Param("id")
// 	if postID == "" {
// 					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 					return
// 	}
// 	postIDUint, err := strconv.Atoi(postID)

// 	comments, total, err := h.repo.GetCommentsByPostID(uint(postIDUint), 1, 10)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"data": comments,
// 		"meta": gin.H{
// 			"total": total,
// 		},
// 		"page": 1,
// 		"page_size": 10,
// 	})

// }

// // Like 点赞模型
// // type Like struct {
// //     ID        uint      `gorm:"primarykey" json:"id"`
// //     UserID    uint      `json:"user_id"`
// //     TargetID  uint      `json:"target_id"`   // 点赞目标ID（可能是帖子或评论）
// //     Type      int       `json:"type"`        // 1:帖子点赞 2:评论点赞
// //     CreatedAt time.Time `json:"created_at"`
// // }

// // CreateLike 创建点赞
// func (h *Handler) CreateLike(c *gin.Context) {
// 	postID := c.Param("id")
// 	if postID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 		return
// 	}

// 	userID := c.GetUint("userID")
// 	postIDUint, err := strconv.Atoi(postID)

// 	like := &model.Like{
// 		UserID:   userID,
// 		TargetID: uint(postIDUint),
// 		Type:     1, // 假设这里是帖子点赞，如果是评论则为2
// 	}

// 	err = h.repo.CreateLike(like)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create like"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Like created successfully"})
// }

// // DeleteLike 删除点赞
// func (h *Handler) DeleteLike(c *gin.Context) {
// 	postID := c.Param("id")
// 	if postID == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
// 		return
// 	}

// 	userID := c.GetUint("userID")

// 	postIDUint, err := strconv.Atoi(postID)
// 	err = h.repo.DeleteLike(uint(postIDUint), userID, 1)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete like"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
// }
