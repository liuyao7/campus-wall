// // internal/handler/topic.go

package handler

// import (
// 	"campus-wall/internal/middleware"
// 	"campus-wall/internal/model"
// 	"campus-wall/internal/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// type TopicHandler struct {
//     topicService *service.TopicService
// }

// func NewTopicHandler(topicService *service.TopicService) *TopicHandler {
//     return &TopicHandler{
//         topicService: topicService,
//     }
// }

// // CreateTopic 创建话题
// func (h *TopicHandler) CreateTopic(c *gin.Context) {
//     userID, _ := c.Get(middleware.UserIDKey)

//     var req model.CreateTopicRequest
//     if err := c.ShouldBindJSON(&req); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     topic, err := h.topicService.CreateTopic(userID.(uint), &req)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusCreated, topic)
// }

// // GetTopicList 获取话题列表
// func (h *TopicHandler) GetTopicList(c *gin.Context) {
//     var query model.TopicListQuery
//     if err := c.ShouldBindQuery(&query); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     topics, total, err := h.topicService.GetTopicList(&query)
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{
//         "topics": topics,
//         "total":  total,
//     })
// }

// // FollowTopic 关注话题
// func (h *TopicHandler) FollowTopic(c *gin.Context) {
//     userID, _ := c.Get(middleware.UserIDKey)
//     topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
//         return
//     }

//     if err := h.topicService.FollowTopic(userID.(uint), uint(topicID)); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "success"})
// }

// // UnfollowTopic 取消关注话题
// func (h *TopicHandler) UnfollowTopic(c *gin.Context) {
//     userID, _ := c.Get(middleware.UserIDKey)
//     topicID, err := strconv.ParseUint(c.Param("id"), 10, 64)
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": "invalid topic id"})
//         return
//     }

//     if err := h.topicService.UnfollowTopic(userID.(uint), uint(topicID)); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, gin.H{"message": "success"})
// }

// // GetUserFollowedTopics 获取用户关注的话题
// func (h *TopicHandler) GetUserFollowedTopics(c *gin.Context) {
//     userID, _ := c.Get(middleware.UserIDKey)

//     topics, err := h.topicService.GetUserFollowedTopics(userID.(uint))
//     if err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, topics)
// }