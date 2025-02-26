// // internal/service/comment.go

package service

// import (
// 	"campus-wall/internal/model"
// 	"campus-wall/pkg/errors"
// 	"context"
// 	"fmt"

// 	"github.com/go-redis/redis/v8"
// 	"gorm.io/gorm"
// )

// type CommentService struct {
//     db    *gorm.DB
//     cache *redis.Client
// }

// func NewCommentService(db *gorm.DB, cache *redis.Client) *CommentService {
//     return &CommentService{
//         db:    db,
//         cache: cache,
//     }
// }

// // CreateComment 创建评论
// func (s *CommentService) CreateComment(userID uint, req *model.CommentRequest, postID uint) (*model.Comment, error) {
//     // 1. 检查帖子是否存在
//     var post model.Post
//     if err := s.db.First(&post, postID).Error; err != nil {
//         if err == gorm.ErrRecordNotFound {
//             return nil, errors.NewNotFoundError("post not found")
//         }
//         return nil, err
//     }

//     // 2. 如果是回复，检查父评论是否存在
//     if req.ParentID != nil {
//         var parent model.Comment
//         if err := s.db.First(&parent, *req.ParentID).Error; err != nil {
//             if err == gorm.ErrRecordNotFound {
//                 return nil, errors.NewNotFoundError("parent comment not found")
//             }
//             return nil, err
//         }
//         // 只允许二级评论
//         if parent.ParentID != nil {
//             return nil, errors.NewBadRequestError("can only reply to first level comment")
//         }
//     }

//     // 3. 创建评论
//     comment := &model.Comment{
//         PostID:   postID,
//         UserID:   userID,
//         ParentID: req.ParentID,
//         Content:  req.Content,
//     }

//     err := s.db.Transaction(func(tx *gorm.DB) error {
//         // 创建评论
//         if err := tx.Create(comment).Error; err != nil {
//             return err
//         }

//         // 更新帖子评论数
//         if err := tx.Model(&post).UpdateColumn("comments", gorm.Expr("comments + ?", 1)).Error; err != nil {
//             return err
//         }

//         return nil
//     })

//     if err != nil {
//         return nil, err
//     }

//     // 4. 加载关联数据
//     if err := s.db.Preload("User").First(comment, comment.ID).Error; err != nil {
//         return nil, err
//     }

//     // 5. 清除缓存
//     s.clearCommentCache(postID)

//     return comment, nil
// }

// // DeleteComment 删除评论
// func (s *CommentService) DeleteComment(userID, commentID uint) error {
//     var comment model.Comment
//     if err := s.db.First(&comment, commentID).Error; err != nil {
//         if err == gorm.ErrRecordNotFound {
//             return errors.NewNotFoundError("comment not found")
//         }
//         return err
//     }

//     // 检查权限
//     if comment.UserID != userID {
//         return errors.NewForbiddenError("no permission to delete this comment")
//     }

//     return s.db.Transaction(func(tx *gorm.DB) error {
//         // 软删除评论
//         if err := tx.Model(&comment).Update("status", 2).Error; err != nil {
//             return err
//         }

//         // 更新帖子评论数
//         if err := tx.Model(&model.Post{}).
//             Where("id = ?", comment.PostID).
//             UpdateColumn("comments", gorm.Expr("comments - ?", 1)).
//             Error; err != nil {
//             return err
//         }

//         // 清除缓存
//         s.clearCommentCache(comment.PostID)

//         return nil
//     })
// }

// // GetCommentList 获取评论列表
// func (s *CommentService) GetCommentList(query *model.CommentListQuery, currentUserID *uint) (*model.CommentListResponse, error) {
//     // 1. 尝试从缓存获取
//     cacheKey := s.generateCommentCacheKey(query)
//     if response, err := s.getCommentsFromCache(cacheKey); err == nil {
//         return response, nil
//     }

//     // 2. 构建查询
//     db := s.db.Model(&model.Comment{}).
//         Preload("User").
//         Where("post_id = ? AND status = 1", query.PostID)

//     // 3. 只查询一级评论
//     db = db.Where("parent_id IS NULL")

//     // 4. 排序
//     switch query.Sort {
//     case "hot":
//         db = db.Order("likes DESC")
//     default:
//         db = db.Order("created_at DESC")
//     }

//     // 5. 获取总数
//     var total int64
//     if err := db.Count(&total).Error; err != nil {
//         return nil, err
//     }

//     // 6. 分页查询
//     var comments []model.Comment
//     if err := db.Offset((query.Page - 1) * query.Size).
//         Limit(query.Size).
//         Find(&comments).Error; err != nil {
//         return nil, err
//     }

//     // 7. 加载回复
//     if err := s.loadReplies(&comments); err != nil {
//         return nil, err
//     }

//     // 8. 处理点赞状态
//     if currentUserID != nil {
//         if err := s.processLikeStatus(&comments, *currentUserID); err != nil {
//             return nil, err
//         }
//     }

//     response := &model.CommentListResponse{
//         Comments: comments,
//         Total:   int(total),
//     }
//     // 9. 存入缓存
//     go s.setCommentsCache(cacheKey, response)

//     return response, nil
// }

// // LikeComment 点赞评论
// func (s *CommentService) LikeComment(userID, commentID uint) error {
//     lockKey := fmt.Sprintf("lock:comment:like:%d:%d", userID, commentID)
//     if !s.acquireLock(lockKey) {
//         return errors.NewBadRequestError("operation too frequent")
//     }
//     defer s.releaseLock(lockKey)

//     return s.db.Transaction(func(tx *gorm.DB) error {
//         // 检查是否已点赞
//         var exists bool
//         if err := tx.Model(&model.CommentLike{}).
//             Where("user_id = ? AND comment_id = ?", userID, commentID).
//             Select("1").
//             Scan(&exists).Error; err != nil {
//             return err
//         }

//         if exists {
//             return errors.NewBadRequestError("already liked")
//         }

//         // 创建点赞记录
//         like := &model.CommentLike{
//             CommentID: commentID,
//             UserID:    userID,
//         }
//         if err := tx.Create(like).Error; err != nil {
//             return err
//         }

//         // 更新评论点赞数
//         if err := tx.Model(&model.Comment{}).
//             Where("id = ?", commentID).
//             UpdateColumn("likes", gorm.Expr("likes + ?", 1)).
//             Error; err != nil {
//             return err
//         }

//         return nil
//     })
// }

// // UnlikeComment 取消点赞评论
// func (s *CommentService) UnlikeComment(userID, commentID uint) error {
//     // ... 实现类似LikeComment的逻辑
// }

// // 辅助方法
// func (s *CommentService) loadReplies(comments *[]model.Comment) error {
//     for i := range *comments {
//         if err := s.db.Preload("User").
//             Where("parent_id = ? AND status = 1", (*comments)[i].ID).
//             Order("created_at ASC").
//             Find(&(*comments)[i].Replies).Error; err != nil {
//             return err
//         }
//     }
//     return nil
// }

// func (s *CommentService) processLikeStatus(comments *[]model.Comment, userID uint) error {
//     // 获取用户对评论的点赞状态
//     commentIDs := make([]uint, 0)
//     for _, comment := range *comments {
//         commentIDs = append(commentIDs, comment.ID)
//         for _, reply := range comment.Replies {
//             commentIDs = append(commentIDs, reply.ID)
//         }
//     }

//     var likes []model.CommentLike
//     if err := s.db.Where("user_id = ? AND comment_id IN ?", userID, commentIDs).
//         Find(&likes).Error; err != nil {
//         return err
//     }

//     likeMap := make(map[uint]bool)
//     for _, like := range likes {
//         likeMap[like.CommentID] = true
//     }

//     // 设置点赞状态
//     for i := range *comments {
//         (*comments)[i].IsLiked = likeMap[(*comments)[i].ID]
//         for j := range (*comments)[i].Replies {
//             (*comments)[i].Replies[j].IsLiked = likeMap[(*comments)[i].Replies[j].ID]
//         }
//     }

//     return nil
// }

// // 缓存相关方法
// func (s *CommentService) generateCommentCacheKey(query *model.CommentListQuery) string {
//     return fmt.Sprintf("comments:%d:%s:%d:%d", query.PostID, query.Sort, query.Page, query.Size)
// }

// func (s *CommentService) clearCommentCache(postID uint) {
//     pattern := fmt.Sprintf("comments:%d:*", postID)
//     keys, _ := s.cache.Keys(context.Background(), pattern).Result()
//     if len(keys) > 0 {
//         s.cache.Del(context.Background(), keys...)
//     }
// }

// func (s *CommentService) getCommentsFromCache(key string) (*model.CommentListResponse, error) {
//     // 从缓存中获取数据的实现
//     return nil, nil
// }

// func (s *CommentService) setCommentsCache(key string, response *model.CommentListResponse) {
//     // 将数据存储到缓存中的实现
//     return
// }

// func (s *CommentService) acquireLock(key string) bool {
//     // 锁的获取逻辑
//     return true
// }

// func (s *CommentService) releaseLock(key string) {
//     // 锁的释放逻辑
//     return
// }