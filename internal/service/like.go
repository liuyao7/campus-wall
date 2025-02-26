// // internal/service/like.go

package service

// import (
// 	"campus-wall/internal/model"
// 	"campus-wall/pkg/errors"
// 	"context"
// 	"fmt"
// 	"time"

// 	"github.com/go-redis/redis/v8"
// 	"gorm.io/gorm"
// )

// type LikeService struct {
//     db    *gorm.DB
//     cache *redis.Client
// }

// func NewLikeService(db *gorm.DB, cache *redis.Client) *LikeService {
//     return &LikeService{
//         db:    db,
//         cache: cache,
//     }
// }

// // LikePost 点赞帖子
// func (s *LikeService) LikePost(userID, postID uint) error {
//     // 1. 检查帖子是否存在
//     var post model.Post
//     if err := s.db.First(&post, postID).Error; err != nil {
//         if err == gorm.ErrRecordNotFound {
//             return errors.NewNotFoundError("post not found")
//         }
//         return err
//     }

//     // 2. 使用Redis作为分布式锁，防止重复点赞
//     lockKey := fmt.Sprintf("lock:like:%d:%d", userID, postID)
//     if !s.acquireLock(lockKey) {
//         return errors.NewBadRequestError("operation too frequent")
//     }
//     defer s.releaseLock(lockKey)

//     // 3. 开启事务
//     return s.db.Transaction(func(tx *gorm.DB) error {
//         // 检查是否已点赞
//         var exists bool
//         if err := tx.Model(&model.PostLike{}).
//             Where("user_id = ? AND post_id = ?", userID, postID).
//             Select("1").
//             Scan(&exists).Error; err != nil {
//             return err
//         }

//         if exists {
//             return errors.NewBadRequestError("already liked")
//         }

//         // 创建点赞记录
//         like := &model.PostLike{
//             PostID: postID,
//             UserID: userID,
//         }
//         if err := tx.Create(like).Error; err != nil {
//             return err
//         }

//         // 更新帖子点赞数
//         if err := tx.Model(&post).UpdateColumn("likes", gorm.Expr("likes + ?", 1)).Error; err != nil {
//             return err
//         }

//         // 更新缓存
//         s.updateLikeCache(postID, userID, true)

//         return nil
//     })
// }

// // UnlikePost 取消点赞
// func (s *LikeService) UnlikePost(userID, postID uint) error {
//     // 使用分布式锁
//     lockKey := fmt.Sprintf("lock:unlike:%d:%d", userID, postID)
//     if !s.acquireLock(lockKey) {
//         return errors.NewBadRequestError("operation too frequent")
//     }
//     defer s.releaseLock(lockKey)

//     return s.db.Transaction(func(tx *gorm.DB) error {
//         // 删除点赞记录
//         result := tx.Where("user_id = ? AND post_id = ?", userID, postID).
//             Delete(&model.PostLike{})

//         if result.Error != nil {
//             return result.Error
//         }

//         if result.RowsAffected == 0 {
//             return errors.NewBadRequestError("not liked yet")
//         }

//         // 更新帖子点赞数
//         if err := tx.Model(&model.Post{}).
//             Where("id = ?", postID).
//             UpdateColumn("likes", gorm.Expr("likes - ?", 1)).
//             Error; err != nil {
//             return err
//         }

//         // 更新缓存
//         s.updateLikeCache(postID, userID, false)

//         return nil
//     })
// }

// // IsLiked 检查用户是否已点赞
// func (s *LikeService) IsLiked(userID, postID uint) (bool, error) {
//     // 1. 先查询缓存
//     liked, err := s.getLikeFromCache(postID, userID)
//     if err == nil {
//         return liked, nil
//     }

//     // 2. 缓存未命中，查询数据库
//     var exists bool
//     err = s.db.Model(&model.PostLike{}).
//         Where("user_id = ? AND post_id = ?", userID, postID).
//         Select("1").
//         Scan(&exists).Error

//     if err != nil {
//         return false, err
//     }

//     // 3. 更新缓存
//     s.updateLikeCache(postID, userID, exists)

//     return exists, nil
// }

// // GetPostLikeUsers 获取帖子点赞用户列表
// func (s *LikeService) GetPostLikeUsers(postID uint, page, size int) ([]model.User, int64, error) {
//     var total int64
//     if err := s.db.Model(&model.PostLike{}).
//         Where("post_id = ?", postID).
//         Count(&total).Error; err != nil {
//         return nil, 0, err
//     }

//     var users []model.User
//     err := s.db.Model(&model.User{}).
//         Joins("JOIN post_likes ON users.id = post_likes.user_id").
//         Where("post_likes.post_id = ?", postID).
//         Offset((page - 1) * size).
//         Limit(size).
//         Find(&users).Error

//     return users, total, err
// }

// // GetUserLikedPosts 获取用户点赞的帖子列表
// func (s *LikeService) GetUserLikedPosts(userID uint, page, size int) ([]model.Post, int64, error) {
//     var total int64
//     if err := s.db.Model(&model.PostLike{}).
//         Where("user_id = ?", userID).
//         Count(&total).Error; err != nil {
//         return nil, 0, err
//     }

//     var posts []model.Post
//     err := s.db.Model(&model.Post{}).
//         Preload("User").
//         Preload("Images").
//         Joins("JOIN post_likes ON posts.id = post_likes.post_id").
//         Where("post_likes.user_id = ?", userID).
//         Offset((page - 1) * size).
//         Limit(size).
//         Find(&posts).Error

//     return posts, total, err
// }

// // Redis相关辅助方法
// func (s *LikeService) acquireLock(key string) bool {
//     return s.cache.SetNX(context.Background(), key, 1, time.Second*5).Val()
// }

// func (s *LikeService) releaseLock(key string) {
//     s.cache.Del(context.Background(), key)
// }

// func (s *LikeService) getLikeFromCache(postID, userID uint) (bool, error) {
//     key := fmt.Sprintf("like:%d:%d", postID, userID)
//     val, err := s.cache.Get(context.Background(), key).Int()
//     if err != nil {
//         return false, err
//     }
//     return val == 1, nil
// }

// func (s *LikeService) updateLikeCache(postID, userID uint, liked bool) {
//     key := fmt.Sprintf("like:%d:%d", postID, userID)
//     val := 0
//     if liked {
//         val = 1
//     }
//     s.cache.Set(context.Background(), key, val, time.Hour*24)
// }