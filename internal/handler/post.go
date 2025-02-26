// // internal/handler/post.go

package handler

// import (
// 	"campus-wall/internal/middleware"
// 	"campus-wall/internal/model"
// 	"campus-wall/internal/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type PostHandler struct {
//     postService *service.PostService
// }

// func NewPostHandler(postService *service.PostService) *PostHandler {
//     return &PostHandler{
//         postService: postService,
//     }
// }

// // UploadImages 上传帖子图片
// func (h *PostHandler) UploadImages(c *gin.Context) {
//     // 获取上传的文件
//     form, err := c.MultipartForm()
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
//         return
//     }

//     files := form.File["images"]
//     if len(files) == 0 {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
//         return
//     }

//     // 限制最大上传数量
//     if len(files) > 9 {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "too many files"})
//         return
//     }

//     // 上传图片
//     urls, err := h.postService.UploadPostImages(files)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"urls": urls})
// }

// // CreatePost 创建帖子
// func (h *PostHandler) CreatePost(c *gin.Context) {
//     // 1. 获取当前用户ID
//     userID, exists := c.Get(middleware.UserIDKey)
//     if !exists {
//         c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
//         return
//     }

//     // 2. 验证请求数据
//     var req model.CreatePostRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // 3. 创建帖子
//     post, err := h.postService.CreatePost(userID.(uint), &req)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusCreated, post)
// }

// // GetPostList 获取帖子列表
// func (h *PostHandler) GetPostList(c *gin.Context) {
//     var query model.PostListQuery
//     if err := c.ShouldBindQuery(&query); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     response, err := h.postService.GetPostList(&query)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, response)
// }

// // GetUserPosts 获取用户的帖子列表
// func (h *PostHandler) GetUserPosts(c *gin.Context) {
//     userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
//         return
//     }

//     var query model.PostListQuery
//     if err := c.ShouldBindQuery(&query); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     response, err := h.postService.GetUserPosts(uint(userID), &query)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, response)
// }

// // GetTopicPosts 获取话题下的帖子列表
// func (h *PostHandler) GetTopicPosts(c *gin.Context) {
//     topicID, err := strconv.ParseUint(c.Param("topic_id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
//         return
//     }

//     var query model.PostListQuery
//     if err := c.ShouldBindQuery(&query); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     response, err := h.postService.GetTopicPosts(uint(topicID), &query)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, response)
// }